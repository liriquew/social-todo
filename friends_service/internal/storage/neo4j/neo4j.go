package neo_storage

import (
	"context"

	config "github.com/liriquew/social-todo/friends_service/internal/lib/config"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Storage struct {
	driver neo4j.DriverWithContext
	dbName string
}

func New(config config.Neo4jConfig) (*Storage, error) {
	uri := "neo4j://localhost:" + config.Port

	neoDriver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(config.Username, config.Password, ""))
	if err != nil {
		return nil, err
	}

	err = neoDriver.VerifyConnectivity(context.Background())
	if err != nil {
		return nil, err
	}

	return &Storage{driver: neoDriver, dbName: config.DBName}, nil
}

func (s *Storage) Close() error {
	c := context.Background()
	return s.driver.Close(c)
}

func (s *Storage) AddFriend(ctx context.Context, UID, friendID int64) error {
	_, err := neo4j.ExecuteQuery(ctx, s.driver, `
		MERGE (u1:User {id: $UID}) 
		MERGE (u2:User {id: $FID}) 
		MERGE (u1)-[r:FRIEND]->(u2)`,
		map[string]any{
			"UID": UID,
			"FID": friendID,
		}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))

	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) RemoveFriend(ctx context.Context, UID, friendID int64) error {
	_, err := neo4j.ExecuteQuery(ctx, s.driver, `
		MATCH (u1:User{id: $UID})
		MATCH (u1)-[r:FRIEND]-(b:User{id: $FID})
		DELETE r`,
		map[string]any{
			"UID": UID,
			"FID": friendID,
		}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))

	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) ListFriends(ctx context.Context, UID int64) ([]int64, error) {
	resp, err := neo4j.ExecuteQuery(ctx, s.driver, `
		MATCH (u: User{id: $UID}) 
		MATCH (u)-[:FRIEND]-(fr:User)
		return fr.id`,
		map[string]any{
			"UID": UID,
		}, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j"))

	if err != nil {
		return nil, err
	}

	ans := make([]int64, 0, len(resp.Records))
	for _, record := range resp.Records {
		ans = append(ans, record.AsMap()["fr.id"].(int64))
	}
	return ans, nil
}
