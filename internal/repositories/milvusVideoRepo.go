package repositories

import (
	"D0SL_organizer/internal/config"
	"D0SL_organizer/pkg/client/models"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type MilvusVideoRepo struct {
	client *client.Client
}

const collectionName = "BaseCollection"
const vectorDim = 768

func NewMilvusRepo(cfg *config.Config) (VideoRepo, error) {
	milvusAddr := fmt.Sprintf("%s:%s", cfg.Milvus.Host, cfg.Milvus.Port)

	client, err := client.NewClient(context.Background(), client.Config{
		Address: milvusAddr,
	})
	if err != nil {
		return nil, err
	}

	fmt.Println("client created")

	InitMilvusClient(context.Background(), &client)

	return &MilvusVideoRepo{client: &client}, nil
}

func InitMilvusClient(ctx context.Context, client *client.Client) *client.Client {
	has, err := (*client).HasCollection(ctx, collectionName)
	if err != nil {
		log.Fatal("failed to check whether collection exists:", err.Error())
	}
	if has {
		// collection with same name exist, clean up mess
		fmt.Println("already exist")
		// _ = (*client).DropCollection(ctx, collectionName)
		return client
	}

	// define collection schema, see film.csv
	schema := entity.NewSchema().WithName(collectionName).WithDescription("this is the example collection for insert and search").
		WithField(entity.NewField().WithName("ID").WithDataType(entity.FieldTypeInt64).WithIsPrimaryKey(true)).
		WithField(entity.NewField().WithName("Link").WithDataType(entity.FieldTypeVarChar).WithMaxLength(2000)).
		WithField(entity.NewField().WithName("Description").WithDataType(entity.FieldTypeVarChar).WithMaxLength(2000)).
		WithField(entity.NewField().WithName("Vector").WithDataType(entity.FieldTypeFloatVector).WithDim(vectorDim))

	err = (*client).CreateCollection(ctx, schema, entity.DefaultShardNumber) // only 1 shard
	if err != nil {
		log.Fatal("failed to create collection:", err.Error())
	}

	videos, err := loadVideoCSV()
	if err != nil {
		log.Fatal("failed to load film data csv:", err.Error())
	}

	fmt.Println("count: ", len(videos))
	// row-base covert to column-base
	ids := make([]int64, 0, len(videos))
	links := make([]string, 0, len(videos))
	descriptions := make([]string, 0, len(videos))
	vectors := make([][]float32, 0, len(videos))
	for idx, video := range videos {
		ids = append(ids, video.ID)
		links = append(links, video.Link)
		// idTitle[film.ID] = film.Title
		descriptions = append(descriptions, video.Description)
		vectors = append(vectors, videos[idx].Vector[:]) // prevent same vector
	}
	idColumn := entity.NewColumnInt64("ID", ids)
	linkColumn := entity.NewColumnVarChar("Link", links)
	descriptionColumn := entity.NewColumnVarChar("Description", descriptions)
	vectorColumn := entity.NewColumnFloatVector("Vector", vectorDim, vectors)

	// insert into default partition
	_, err = (*client).Insert(ctx, collectionName, "", idColumn, linkColumn, descriptionColumn, vectorColumn)
	if err != nil {
		log.Fatal("failed to insert film data:", err.Error())
	}
	log.Println("insert completed")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
	err = (*client).Flush(ctx, collectionName, false)
	if err != nil {
		log.Fatal("failed to flush collection:", err.Error())
	}
	log.Println("flush completed")

	count, err := (*client).GetCollectionStatistics(ctx, collectionName)
	if err != nil {
		log.Fatal("failed to count entities:", err.Error())
	}
	fmt.Println("Total number of entities in the collection:", count)

	return client
}

func loadVideoCSV() ([]models.Video, error) {
	f, err := os.Open("db_vector.csv")
	if err != nil {
		return []models.Video{}, err
	}

	r := csv.NewReader(f)
	raw, err := r.ReadAll()
	if err != nil {
		return []models.Video{}, err
	}

	videos := make([]models.Video, 0, len(raw))

	for _, line := range raw {

		if len(line) < 2 { // insuffcient column
			continue
		}
		fi := models.Video{}

		// ID
		v, err := strconv.ParseInt(line[0], 10, 64)
		if err != nil {
			continue
		}
		fi.ID = v

		fi.Link = line[1]
		fi.Description = line[2]

		// Vector
		vectorStr := strings.ReplaceAll(line[3], "[", "")
		vectorStr = strings.ReplaceAll(vectorStr, "]", "")
		parts := strings.Split(vectorStr, ",")
		if len(parts) != vectorDim { // dim must be eq vectorDim
			continue
		}
		for idx, part := range parts {
			part = strings.TrimSpace(part)
			v, err := strconv.ParseFloat(part, 32)
			if err != nil {
				continue
			}
			fi.Vector[idx] = float32(v)
		}
		// fmt.Println("Video: ", fi)
		videos = append(videos, fi)
	}
	return videos, nil
}

func (m *MilvusVideoRepo) AddVideo(video models.Video) error {
	ids := []int64{video.ID}
	links := []string{video.Link}
	descriptions := []string{video.Description}
	vectors := [][]float32{video.Vector[:]}

	idColumn := entity.NewColumnInt64("ID", ids)
	linkColumn := entity.NewColumnVarChar("Link", links)
	descriptionColumn := entity.NewColumnVarChar("Description", descriptions)
	vectorColumn := entity.NewColumnFloatVector("Vector", vectorDim, vectors)

	// insert into default partition
	_, err := (*m.client).Insert(context.Background(), collectionName, "", idColumn, linkColumn, descriptionColumn, vectorColumn)
	if err != nil {
		log.Fatal("failed to insert new video:", err.Error())
	}
	log.Println("insert completed")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
	err = (*m.client).Flush(ctx, collectionName, false)
	if err != nil {
		log.Fatal("failed to flush new collection:", err.Error())
	}
	log.Println("flush1 completed")

	count, err := (*m.client).GetCollectionStatistics(ctx, collectionName)
	if err != nil {
		log.Fatal("failed to count entities:", err.Error())
	}
	fmt.Println("Total number of entities in the collection:", count)
	return nil
}

func (m *MilvusVideoRepo) GetSimilarVideosByVector(embedding []float32) ([]string, error) {
	idx, err := entity.NewIndexIvfFlat(
		entity.COSINE,
		1024,
	)
	if err != nil {
		fmt.Println(err)
	}

	err = (*m.client).CreateIndex(context.Background(), collectionName, "Vector", idx, false)
	if err != nil {
		log.Fatal("failed to create index:", err.Error())
	}
	log.Println("index creation completed")

	ctx := context.Background()
	err = (*m.client).LoadCollection(context.Background(), collectionName, false)
	if err != nil {
		return nil, err
	}
	log.Println("load collection completed")

	vector := entity.FloatVector(embedding)
	sp, _ := entity.NewIndexFlatSearchParam()
	sr, err := (*m.client).Search(ctx, collectionName, []string{}, "", []string{"Link"}, []entity.Vector{vector}, "Vector",
		entity.COSINE, 5, sp)
	if err != nil {
		return nil, err
	}

	linksArray := make([]string, 0, 5)
	for _, result := range sr {
		var links *entity.ColumnVarChar
		for _, field := range result.Fields {
			fmt.Println(field.Name())
			if field.Name() == "Link" {
				c, ok := field.(*entity.ColumnVarChar)
				if ok {
					links = c
				}
			}
		}
		if links == nil {
			fmt.Println("No")
			continue
		}
		for i := 0; i < result.ResultCount; i++ {
			link, err := links.ValueByIdx(i)
			if err != nil {
				return nil, err
			}
			linksArray = append(linksArray, link)
		}
	}

	fmt.Println("result: ", linksArray)
	return linksArray, nil
}
