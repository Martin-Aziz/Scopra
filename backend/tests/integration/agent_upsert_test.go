package integration
package integration_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/martin-aziz/scopra/backend/src/models"
	"github.com/martin-aziz/scopra/backend/src/repositories"
)

func TestUpsertAgent_WithUUID(t *testing.T) {






















































































}	})		assert.Equal(t, 5, count, "should have 5 active agents after multiple upserts")		require.NoError(t, err)		err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM agents WHERE status = 'active'").Scan(&count)		var count int		// Verify all 5 agents exist		}			assert.NoError(t, err, "UpsertAgent %d should succeed without type errors", i+1)			err := repo.UpsertAgent(ctx, agentID)			agentID := uuid.New().String()		for i := 0; i < 5; i++ {	t.Run("multiple_upserts_no_type_errors", func(t *testing.T) {	// Test case 3: Multiple agents can be upserted without type errors	})		assert.Equal(t, "active", status, "status should be reset to 'active' on upsert conflict")		require.NoError(t, err, "agent should exist")		).Scan(&status)			agentID,			"SELECT status FROM agents WHERE id = $1::uuid",			ctx,		err = conn.QueryRow(		var status string		assert.NoError(t, err, "second UpsertAgent should succeed (upsert conflict)")		err = repo.UpsertAgent(ctx, agentID)		// Upsert should set status back to active		require.NoError(t, err, "manual status update should succeed")		)			agentID,			"UPDATE agents SET status = 'revoked' WHERE id = $1::uuid",			ctx,		_, err = conn.Exec(		// Manually set status to revoked to verify upsert resets it		require.NoError(t, err, "first UpsertAgent should succeed")		err := repo.UpsertAgent(ctx, agentID)		// First insert		agentID := uuid.New().String()	t.Run("upsert_existing_agent_updates_status", func(t *testing.T) {	// Test case 2: Upsert existing agent should update status without error	})		assert.True(t, len(name) > 6, "agent name should have agent- prefix + uuid substring")		// Name should follow pattern: agent-{8-char-uuid-prefix}		assert.NotEmpty(t, name, "agent name should not be empty")		assert.Equal(t, "active", status, "agent status should be 'active'")		require.NoError(t, err, "agent should exist after insert")		).Scan(&name, &status)			agentID,			"SELECT name, status FROM agents WHERE id = $1::uuid",			ctx,		err = conn.QueryRow(		var name, status string		// Verify agent was inserted with correct name pattern		assert.NoError(t, err, "UpsertAgent should succeed with valid UUID")		err := repo.UpsertAgent(ctx, agentID)		agentID := uuid.New().String()	t.Run("insert_new_agent_with_uuid_substring", func(t *testing.T) {	// Test case 1: Insert new agent with valid UUID	repo := repositories.NewPostgresRepository(conn)	// Create repository	require.NoError(t, err, "agents table does not exist or query failed")	err = conn.QueryRow(ctx, "SELECT COUNT(*) FROM agents").Scan(&count)	var count int	// Verify the agents table exists and is empty	defer conn.Close(ctx)	require.NoError(t, err, "failed to connect to Postgres")	conn, err := pgx.Connect(ctx, dbURL)	// Connect to Postgres	ctx := context.Background()	}		t.Skip("DATABASE_URL not set; skipping Postgres integration test")	if dbURL == "" {	dbURL := os.Getenv("DATABASE_URL")	// Skip if DATABASE_URL is not set (e.g., in CI without Postgres)