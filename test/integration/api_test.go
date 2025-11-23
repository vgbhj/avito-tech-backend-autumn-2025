package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/avito-tech-backend-autumn-2025/test/helpers"
)

func TestAPI_TeamEndpoints(t *testing.T) {
	db, cleanup, err := helpers.SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	router := helpers.SetupTestApp(db)

	// Тест проверяет успешное создание команды с участниками.
	// Ожидается: команда создана, все участники добавлены, возвращается статус 201.
	t.Run("CreateTeam - success", func(t *testing.T) {
		helpers.CleanupDB(db)

		reqBody := map[string]interface{}{
			"team_name": "backend",
			"members": []map[string]interface{}{
				{"user_id": "u1", "username": "Alice", "is_active": true},
				{"user_id": "u2", "username": "Bob", "is_active": true},
			},
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.NotNil(t, response["team"])
		team := response["team"].(map[string]interface{})
		assert.Equal(t, "backend", team["team_name"])
		assert.Len(t, team["members"], 2)
	})

	// Тест проверяет обработку попытки создать команду с уже существующим именем.
	// Ожидается: возвращается ошибка TEAM_EXISTS со статусом 400.
	t.Run("CreateTeam - duplicate team", func(t *testing.T) {
		helpers.CleanupDB(db)

		reqBody := map[string]interface{}{
			"team_name": "backend",
			"members": []map[string]interface{}{
				{"user_id": "u1", "username": "Alice", "is_active": true},
			},
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		req2 := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(body))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusBadRequest, w2.Code)
		var errorResp map[string]interface{}
		json.Unmarshal(w2.Body.Bytes(), &errorResp)
		assert.Equal(t, "TEAM_EXISTS", errorResp["error"].(map[string]interface{})["code"])
	})

	// Тест проверяет успешное получение команды по имени.
	// Ожидается: команда найдена, возвращается со всеми участниками, статус 200.
	t.Run("GetTeam - success", func(t *testing.T) {
		helpers.CleanupDB(db)

		reqBody := map[string]interface{}{
			"team_name": "frontend",
			"members": []map[string]interface{}{
				{"user_id": "u3", "username": "Charlie", "is_active": true},
			},
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		req2 := httptest.NewRequest(http.MethodGet, "/team/get?team_name=frontend", nil)
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusOK, w2.Code)
		var team map[string]interface{}
		json.Unmarshal(w2.Body.Bytes(), &team)
		assert.Equal(t, "frontend", team["team_name"])
	})

	// Тест проверяет обработку запроса несуществующей команды.
	// Ожидается: возвращается ошибка NOT_FOUND со статусом 404.
	t.Run("GetTeam - not found", func(t *testing.T) {
		helpers.CleanupDB(db)

		req := httptest.NewRequest(http.MethodGet, "/team/get?team_name=nonexistent", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestAPI_UserEndpoints(t *testing.T) {
	db, cleanup, err := helpers.SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	router := helpers.SetupTestApp(db)

	// Тест проверяет успешное изменение флага активности пользователя.
	// Ожидается: флаг активности обновлен, возвращается обновленный пользователь, статус 200.
	t.Run("SetActive - success", func(t *testing.T) {
		helpers.CleanupDB(db)

		createTeamReq := map[string]interface{}{
			"team_name": "backend",
			"members": []map[string]interface{}{
				{"user_id": "u1", "username": "Alice", "is_active": true},
			},
		}
		body, _ := json.Marshal(createTeamReq)
		req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		setActiveReq := map[string]interface{}{
			"user_id":   "u1",
			"is_active": false,
		}
		body2, _ := json.Marshal(setActiveReq)
		req2 := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewBuffer(body2))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusOK, w2.Code)
		var response map[string]interface{}
		json.Unmarshal(w2.Body.Bytes(), &response)
		user := response["user"].(map[string]interface{})
		assert.Equal(t, false, user["is_active"])
	})

	// Тест проверяет получение списка PR, где пользователь назначен ревьюером.
	// Ожидается: возвращается список PR пользователя, статус 200.
	t.Run("GetReviews - success", func(t *testing.T) {
		helpers.CleanupDB(db)

		createTeamReq := map[string]interface{}{
			"team_name": "backend",
			"members": []map[string]interface{}{
				{"user_id": "u1", "username": "Alice", "is_active": true},
				{"user_id": "u2", "username": "Bob", "is_active": true},
			},
		}
		body, _ := json.Marshal(createTeamReq)
		req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		createPRReq := map[string]interface{}{
			"pull_request_id":   "pr-1",
			"pull_request_name": "Test PR",
			"author_id":         "u1",
		}
		body2, _ := json.Marshal(createPRReq)
		req2 := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBuffer(body2))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusCreated, w2.Code)

		req3 := httptest.NewRequest(http.MethodGet, "/users/getReview?user_id=u2", nil)
		w3 := httptest.NewRecorder()
		router.ServeHTTP(w3, req3)

		assert.Equal(t, http.StatusOK, w3.Code)
		var response map[string]interface{}
		json.Unmarshal(w3.Body.Bytes(), &response)
		assert.Equal(t, "u2", response["user_id"])
		pullRequests := response["pull_requests"].([]interface{})
		assert.GreaterOrEqual(t, len(pullRequests), 1)
	})
}

func TestAPI_PREndpoints(t *testing.T) {
	db, cleanup, err := helpers.SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	router := helpers.SetupTestApp(db)

	// Тест проверяет создание PR с автоматическим назначением ревьюеров.
	// Ожидается: PR создан, назначено до 2 ревьюеров из команды автора,
	// автор исключен из списка ревьюеров, статус 201.
	t.Run("CreatePR - success with auto-assigned reviewers", func(t *testing.T) {
		helpers.CleanupDB(db)

		createTeamReq := map[string]interface{}{
			"team_name": "backend",
			"members": []map[string]interface{}{
				{"user_id": "u1", "username": "Alice", "is_active": true},
				{"user_id": "u2", "username": "Bob", "is_active": true},
				{"user_id": "u3", "username": "Charlie", "is_active": true},
			},
		}
		body, _ := json.Marshal(createTeamReq)
		req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		createPRReq := map[string]interface{}{
			"pull_request_id":   "pr-1",
			"pull_request_name": "Test PR",
			"author_id":         "u1",
		}
		body2, _ := json.Marshal(createPRReq)
		req2 := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBuffer(body2))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusCreated, w2.Code)
		var response map[string]interface{}
		json.Unmarshal(w2.Body.Bytes(), &response)
		pr := response["pr"].(map[string]interface{})
		assert.Equal(t, "pr-1", pr["pull_request_id"])
		assert.Equal(t, "OPEN", pr["status"])
		reviewers := pr["assigned_reviewers"].([]interface{})
		assert.LessOrEqual(t, len(reviewers), 2)
		assert.GreaterOrEqual(t, len(reviewers), 1)
		// Проверяем, что автор не в списке ревьюеров
		for _, reviewer := range reviewers {
			assert.NotEqual(t, "u1", reviewer)
		}
	})

	// Тест проверяет обработку попытки создать PR с уже существующим ID.
	// Ожидается: возвращается ошибка PR_EXISTS со статусом 409.
	t.Run("CreatePR - duplicate PR", func(t *testing.T) {
		helpers.CleanupDB(db)

		createTeamReq := map[string]interface{}{
			"team_name": "backend",
			"members": []map[string]interface{}{
				{"user_id": "u1", "username": "Alice", "is_active": true},
			},
		}
		body, _ := json.Marshal(createTeamReq)
		req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		createPRReq := map[string]interface{}{
			"pull_request_id":   "pr-1",
			"pull_request_name": "Test PR",
			"author_id":         "u1",
		}
		body2, _ := json.Marshal(createPRReq)
		req2 := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBuffer(body2))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		assert.Equal(t, http.StatusCreated, w2.Code)

		req3 := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBuffer(body2))
		req3.Header.Set("Content-Type", "application/json")
		w3 := httptest.NewRecorder()
		router.ServeHTTP(w3, req3)

		assert.Equal(t, http.StatusConflict, w3.Code)
		var errorResp map[string]interface{}
		json.Unmarshal(w3.Body.Bytes(), &errorResp)
		assert.Equal(t, "PR_EXISTS", errorResp["error"].(map[string]interface{})["code"])
	})

	// Тест проверяет merge PR и идемпотентность операции.
	// Ожидается: PR помечен как MERGED при первом вызове,
	// повторный вызов не приводит к ошибке и возвращает MERGED статус, статус 200.
	t.Run("MergePR - success and idempotent", func(t *testing.T) {
		helpers.CleanupDB(db)

		createTeamReq := map[string]interface{}{
			"team_name": "backend",
			"members": []map[string]interface{}{
				{"user_id": "u1", "username": "Alice", "is_active": true},
			},
		}
		body, _ := json.Marshal(createTeamReq)
		req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		createPRReq := map[string]interface{}{
			"pull_request_id":   "pr-1",
			"pull_request_name": "Test PR",
			"author_id":         "u1",
		}
		body2, _ := json.Marshal(createPRReq)
		req2 := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBuffer(body2))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		mergeReq := map[string]interface{}{
			"pull_request_id": "pr-1",
		}
		body3, _ := json.Marshal(mergeReq)
		req3 := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewBuffer(body3))
		req3.Header.Set("Content-Type", "application/json")
		w3 := httptest.NewRecorder()
		router.ServeHTTP(w3, req3)

		assert.Equal(t, http.StatusOK, w3.Code)
		var response map[string]interface{}
		json.Unmarshal(w3.Body.Bytes(), &response)
		pr := response["pr"].(map[string]interface{})
		assert.Equal(t, "MERGED", pr["status"])

		req4 := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewBuffer(body3))
		req4.Header.Set("Content-Type", "application/json")
		w4 := httptest.NewRecorder()
		router.ServeHTTP(w4, req4)

		assert.Equal(t, http.StatusOK, w4.Code)
		var response2 map[string]interface{}
		json.Unmarshal(w4.Body.Bytes(), &response2)
		pr2 := response2["pr"].(map[string]interface{})
		assert.Equal(t, "MERGED", pr2["status"])
	})

	// Тест проверяет успешное переназначение ревьюера.
	// Ожидается: один ревьюер заменен на другого из той же команды,
	// возвращается новый ревьюер, статус 200.
	t.Run("ReassignReviewer - success", func(t *testing.T) {
		helpers.CleanupDB(db)

		// Создаем команду с достаточным количеством пользователей для переназначения
		// Нужно минимум 4: автор + 2 ревьюера + 1 для замены
		createTeamReq := map[string]interface{}{
			"team_name": "backend",
			"members": []map[string]interface{}{
				{"user_id": "u1", "username": "Alice", "is_active": true},
				{"user_id": "u2", "username": "Bob", "is_active": true},
				{"user_id": "u3", "username": "Charlie", "is_active": true},
				{"user_id": "u4", "username": "David", "is_active": true},
			},
		}
		body, _ := json.Marshal(createTeamReq)
		req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		createPRReq := map[string]interface{}{
			"pull_request_id":   "pr-1",
			"pull_request_name": "Test PR",
			"author_id":         "u1",
		}
		body2, _ := json.Marshal(createPRReq)
		req2 := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBuffer(body2))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		var prResponse map[string]interface{}
		json.Unmarshal(w2.Body.Bytes(), &prResponse)
		pr := prResponse["pr"].(map[string]interface{})
		reviewers := pr["assigned_reviewers"].([]interface{})
		require.Greater(t, len(reviewers), 0)
		oldReviewer := reviewers[0].(string)

		reassignReq := map[string]interface{}{
			"pull_request_id": "pr-1",
			"old_user_id":     oldReviewer,
		}
		body3, _ := json.Marshal(reassignReq)
		req3 := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewBuffer(body3))
		req3.Header.Set("Content-Type", "application/json")
		w3 := httptest.NewRecorder()
		router.ServeHTTP(w3, req3)

		if w3.Code != http.StatusOK {
			var errorResp map[string]interface{}
			json.Unmarshal(w3.Body.Bytes(), &errorResp)
			t.Logf("Error response: %+v", errorResp)
		}
		assert.Equal(t, http.StatusOK, w3.Code, "Response body: %s", w3.Body.String())
		var response map[string]interface{}
		json.Unmarshal(w3.Body.Bytes(), &response)
		assert.NotEmpty(t, response["replaced_by"])
	})

	// Тест проверяет запрет переназначения ревьюеров для объединенного PR.
	// Ожидается: после merge PR нельзя переназначить ревьюеров,
	// возвращается ошибка PR_MERGED со статусом 409.
	t.Run("ReassignReviewer - cannot reassign merged PR", func(t *testing.T) {
		helpers.CleanupDB(db)

		createTeamReq := map[string]interface{}{
			"team_name": "backend",
			"members": []map[string]interface{}{
				{"user_id": "u1", "username": "Alice", "is_active": true},
				{"user_id": "u2", "username": "Bob", "is_active": true},
			},
		}
		body, _ := json.Marshal(createTeamReq)
		req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		createPRReq := map[string]interface{}{
			"pull_request_id":   "pr-1",
			"pull_request_name": "Test PR",
			"author_id":         "u1",
		}
		body2, _ := json.Marshal(createPRReq)
		req2 := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBuffer(body2))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		mergeReq := map[string]interface{}{
			"pull_request_id": "pr-1",
		}
		body3, _ := json.Marshal(mergeReq)
		req3 := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewBuffer(body3))
		req3.Header.Set("Content-Type", "application/json")
		w3 := httptest.NewRecorder()
		router.ServeHTTP(w3, req3)

		reassignReq := map[string]interface{}{
			"pull_request_id": "pr-1",
			"old_user_id":     "u2",
		}
		body4, _ := json.Marshal(reassignReq)
		req4 := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewBuffer(body4))
		req4.Header.Set("Content-Type", "application/json")
		w4 := httptest.NewRecorder()
		router.ServeHTTP(w4, req4)

		assert.Equal(t, http.StatusConflict, w4.Code)
		var errorResp map[string]interface{}
		json.Unmarshal(w4.Body.Bytes(), &errorResp)
		assert.Equal(t, "PR_MERGED", errorResp["error"].(map[string]interface{})["code"])
	})
}

func TestAPI_HealthEndpoint(t *testing.T) {
	db, cleanup, err := helpers.SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	router := helpers.SetupTestApp(db)

	// Тест проверяет health check endpoint.
	// Ожидается: сервис отвечает "OK", статус 200.
	t.Run("Health check", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "OK", w.Body.String())
	})
}
