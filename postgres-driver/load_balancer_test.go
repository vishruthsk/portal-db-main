package postgresdriver

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/vishruthsk/portal-db/types"
)

func (ts *PGDriverTestSuite) Test_ReadLoadBalancers() {
	tests := []struct {
		name          string
		loadBalancers []*types.LoadBalancer
		err           error
	}{
		{
			name: "Should return all Load Balancers from the database ordered by lb_id",
			loadBalancers: []*types.LoadBalancer{
				{
					ID:                "test_lb_34987u329rfn23f",
					Name:              "vipr_app_123",
					UserID:            "test_user_1dbffbdfeeb225",
					ApplicationIDs:    []string{"test_app_47hfnths73j2se"},
					RequestTimeout:    5_000,
					Gigastake:         true,
					GigastakeRedirect: true,
					StickyOptions: types.StickyOptions{
						Duration:      "60",
						StickyOrigins: []string{"chrome-extension://", "moz-extension://"},
						StickyMax:     300,
						Stickiness:    true,
					},
					Users: []types.UserAccess{
						{RoleName: "OWNER", UserID: "test_user_1dbffbdfeeb225", Email: "owner1@test.com", Accepted: true},
						{RoleName: "ADMIN", UserID: "test_user_admin1234", Email: "admin1@test.com", Accepted: true},
						{RoleName: "MEMBER", UserID: "test_user_member1234", Email: "member1@test.com", Accepted: true},
					},
				},
				{
					ID:                "test_lb_34gg4g43g34g5hh",
					Name:              "test_lb_redirect",
					UserID:            "test_user_redirect233344",
					ApplicationIDs:    []string{""},
					RequestTimeout:    5_000,
					Gigastake:         false,
					GigastakeRedirect: false,
					StickyOptions: types.StickyOptions{
						Duration:      "20",
						StickyOrigins: []string{"test-extension://", "test-extension2://"},
						StickyMax:     600,
						Stickiness:    false,
					},
					Users: []types.UserAccess{
						{RoleName: "OWNER", UserID: "test_user_redirect233344", Email: "owner3@test.com", Accepted: true},
						{RoleName: "MEMBER", UserID: "test_user_member5678", Email: "member2@test.com", Accepted: true},
					},
				},
				{
					ID:                "test_lb_3890ru23jfi32fj",
					Name:              "vipr_app_456",
					UserID:            "test_user_04228205bd261a",
					ApplicationIDs:    []string{"test_app_5hdf7sh23jd828"},
					RequestTimeout:    5_000,
					Gigastake:         true,
					GigastakeRedirect: true,
					StickyOptions: types.StickyOptions{
						Duration:      "40",
						StickyOrigins: []string{"chrome-extension://"},
						StickyMax:     400,
						Stickiness:    true,
					},
					Users: []types.UserAccess{
						{RoleName: "OWNER", UserID: "test_user_04228205bd261a", Email: "owner2@test.com", Accepted: true},
						{RoleName: "ADMIN", UserID: "test_user_admin5678", Email: "admin2@test.com", Accepted: true},
					},
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		loadBalancers, err := ts.driver.ReadLoadBalancers(testCtx)
		ts.Equal(test.err, err)
		for i, loadBalancer := range loadBalancers {
			ts.Equal(test.loadBalancers[i].ID, loadBalancer.ID)
			ts.Equal(test.loadBalancers[i].UserID, loadBalancer.UserID)
			ts.Equal(test.loadBalancers[i].Name, loadBalancer.Name)
			ts.Equal(test.loadBalancers[i].UserID, loadBalancer.UserID)
			ts.Equal(test.loadBalancers[i].ApplicationIDs, loadBalancer.ApplicationIDs)
			ts.Equal(test.loadBalancers[i].RequestTimeout, loadBalancer.RequestTimeout)
			ts.Equal(test.loadBalancers[i].Gigastake, loadBalancer.Gigastake)
			ts.Equal(test.loadBalancers[i].GigastakeRedirect, loadBalancer.GigastakeRedirect)
			ts.Equal(test.loadBalancers[i].StickyOptions, loadBalancer.StickyOptions)
			ts.Equal(test.loadBalancers[i].Users, loadBalancer.Users)
			ts.NotEmpty(loadBalancer.CreatedAt)
			ts.NotEmpty(loadBalancer.UpdatedAt)
		}
	}
}

func (ts *PGDriverTestSuite) Test_ReadUserRoles() {
	tests := []struct {
		name         string
		userRolesMap map[string]map[string][]types.PermissionsEnum
		err          error
	}{
		{
			name: "Should return all Load Balancers from the database ordered by lb_id",
			userRolesMap: map[string]map[string][]types.PermissionsEnum{
				"test_user_1dbffbdfeeb225": {
					"test_lb_34987u329rfn23f": []types.PermissionsEnum{
						types.ReadEndpoint,
						types.WriteEndpoint,
					},
				},
				"test_user_admin1234": {
					"test_lb_34987u329rfn23f": []types.PermissionsEnum{
						types.ReadEndpoint,
						types.WriteEndpoint,
					},
				},
				"test_user_member1234": {
					"test_lb_34987u329rfn23f": []types.PermissionsEnum{
						types.ReadEndpoint},
				},
				"test_user_04228205bd261a": {
					"test_lb_3890ru23jfi32fj": []types.PermissionsEnum{
						types.ReadEndpoint,
						types.WriteEndpoint,
					},
				},
				"test_user_admin5678": {
					"test_lb_3890ru23jfi32fj": []types.PermissionsEnum{
						types.ReadEndpoint,
						types.WriteEndpoint,
					},
				},
				"test_user_redirect233344": {
					"test_lb_34gg4g43g34g5hh": []types.PermissionsEnum{
						types.ReadEndpoint,
						types.WriteEndpoint,
					},
				},
				"test_user_member5678": {
					"test_lb_34gg4g43g34g5hh": []types.PermissionsEnum{
						types.ReadEndpoint},
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		userRoles, err := ts.driver.ReadUserRoles(testCtx)
		ts.Equal(test.err, err)
		ts.Equal(test.userRolesMap, userRoles)
	}
}

func (ts *PGDriverTestSuite) Test_WriteLoadBalancer() {
	tests := []struct {
		name               string
		loadBalancerInputs []*types.LoadBalancer
		expectedNumOfLBs   int
		expectedLB         SelectOneLoadBalancerRow
		err                error
	}{
		{
			name: "Should create a single load balancer successfully with correct input",
			loadBalancerInputs: []*types.LoadBalancer{
				{
					Name:              "vipr_app_789",
					UserID:            "test_user_47fhsd75jd756sh",
					RequestTimeout:    5000,
					Gigastake:         true,
					GigastakeRedirect: true,
					ApplicationIDs:    []string{"test_app_47hfnths73j2se"},
					StickyOptions: types.StickyOptions{
						Duration:      "70",
						StickyOrigins: []string{"chrome-extension://"},
						StickyMax:     400,
						Stickiness:    true,
					},
					Users: []types.UserAccess{
						{
							UserID:   "test_user_47fhsd75jd756sh",
							RoleName: types.RoleOwner,
							Email:    "owner4@test.com",
							Accepted: true,
						},
					},
				},
			},
			expectedNumOfLBs: 4,
			expectedLB: SelectOneLoadBalancerRow{
				Name:              sql.NullString{Valid: true, String: "vipr_app_789"},
				UserID:            sql.NullString{Valid: true, String: "test_user_47fhsd75jd756sh"},
				RequestTimeout:    sql.NullInt32{Valid: true, Int32: 5000},
				Gigastake:         sql.NullBool{Valid: true, Bool: true},
				GigastakeRedirect: sql.NullBool{Valid: true, Bool: true},
				Duration:          sql.NullString{Valid: true, String: "70"},
				StickyMax:         sql.NullInt32{Valid: true, Int32: 400},
				Stickiness:        sql.NullBool{Valid: true, Bool: true},
				Origins:           []string{"chrome-extension://"},
				Users:             json.RawMessage(`[{"email": "owner4@test.com", "userID": "test_user_47fhsd75jd756sh", "accepted": true, "roleName": "OWNER"}]`),
			},
			err: nil,
		},
		{
			name: "Should fail if input does not have at least one user",
			loadBalancerInputs: []*types.LoadBalancer{
				{Users: []types.UserAccess{}},
			},
			err: ErrLBMustHaveUser,
		},
	}

	for _, test := range tests {
		for _, input := range test.loadBalancerInputs {
			createdLB, err := ts.driver.WriteLoadBalancer(testCtx, input)
			if err == nil {
				ts.Equal(test.err, err)
				ts.Len(createdLB.ID, 24)
				ts.Equal(input.Name, createdLB.Name)
				ts.NotEmpty(createdLB.CreatedAt)
				ts.NotEmpty(createdLB.UpdatedAt)

				loadBalancers, err := ts.driver.ReadLoadBalancers(testCtx)
				ts.Equal(test.err, err)
				ts.Len(loadBalancers, test.expectedNumOfLBs)

				loadBalancer, err := ts.driver.SelectOneLoadBalancer(testCtx, createdLB.ID)
				ts.Equal(test.err, err)
				for _, testInput := range test.loadBalancerInputs {
					if testInput.Name == loadBalancer.Name.String {
						ts.Equal(createdLB.ID, loadBalancer.LbID)
						ts.Equal(test.expectedLB.UserID, loadBalancer.UserID)
						ts.Equal(test.expectedLB.Name, loadBalancer.Name)
						ts.Equal(test.expectedLB.UserID, loadBalancer.UserID)
						ts.Equal(test.expectedLB.RequestTimeout, loadBalancer.RequestTimeout)
						ts.Equal(test.expectedLB.Gigastake, loadBalancer.Gigastake)
						ts.Equal(test.expectedLB.GigastakeRedirect, loadBalancer.GigastakeRedirect)
						ts.Equal(test.expectedLB.Duration, loadBalancer.Duration)
						ts.Equal(test.expectedLB.Origins, loadBalancer.Origins)
						ts.Equal(test.expectedLB.StickyMax, loadBalancer.StickyMax)
						ts.Equal(test.expectedLB.Stickiness, loadBalancer.Stickiness)
						ts.Equal(test.expectedLB.Users, loadBalancer.Users)
						ts.NotEmpty(loadBalancer.CreatedAt)
						ts.NotEmpty(loadBalancer.UpdatedAt)
					}
				}
			}
		}
	}
}

func (ts *PGDriverTestSuite) Test_WriteLoadBalancerUser() {
	tests := []struct {
		name              string
		lbIDInput         string
		userInput         types.UserAccess
		expectedUsersJSON json.RawMessage
		expectedUsers     []types.UserAccess
		err               error
	}{
		{
			name:      "Should create a new UserAccess row for a LoadBalancer with correct input",
			lbIDInput: "test_lb_34987u329rfn23f",
			userInput: types.UserAccess{
				UserID:   "test_user_47fhsd75jd756sh",
				RoleName: types.RoleMember,
				Email:    "member5@test.com",
			},
			expectedUsersJSON: json.RawMessage(`[{"email": "owner1@test.com", "userID": "test_user_1dbffbdfeeb225", "accepted": true, "roleName": "OWNER"}, {"email": "admin1@test.com", "userID": "test_user_admin1234", "accepted": true, "roleName": "ADMIN"}, {"email": "member1@test.com", "userID": "test_user_member1234", "accepted": true, "roleName": "MEMBER"}, {"email": "member5@test.com", "userID": "test_user_47fhsd75jd756sh", "accepted": false, "roleName": "MEMBER"}]`),
			expectedUsers: []types.UserAccess{
				{
					UserID:   "test_user_1dbffbdfeeb225",
					RoleName: types.RoleOwner,
					Email:    "owner1@test.com",
					Accepted: true,
				},
				{
					UserID:   "test_user_admin1234",
					RoleName: types.RoleAdmin,
					Email:    "admin1@test.com",
					Accepted: true,
				},
				{
					UserID:   "test_user_member1234",
					RoleName: types.RoleMember,
					Email:    "member1@test.com",
					Accepted: true,
				},
				{
					UserID:   "test_user_47fhsd75jd756sh",
					RoleName: types.RoleMember,
					Email:    "member5@test.com",
					Accepted: false,
				},
			},
			err: nil,
		},
		{
			name:      "Should fail if any input fields are null",
			lbIDInput: "test_lb_34987u329rfn23f",
			userInput: types.UserAccess{
				UserID:   "test_user_47fhsd75jd756sh",
				RoleName: types.RoleMember,
			},
			err: fmt.Errorf("%w: Email", ErrUserInputIsMissingField),
		},
		{
			name:      "Should fail if lb ID not provided",
			lbIDInput: "",
			err:       ErrMissingID,
		},
		{
			name:      "Should fail if attempting to create a User with owner role",
			lbIDInput: "test_lb_3890ru23jfi32fj",
			userInput: types.UserAccess{
				UserID:   "test_user_47fhsd75jd756sh",
				RoleName: types.RoleOwner,
				Email:    "member5@test.com",
			},
			err: ErrCannotSetToOwner,
		},
	}

	for _, test := range tests {
		err := ts.driver.WriteLoadBalancerUser(testCtx, test.lbIDInput, test.userInput)
		ts.Equal(test.err, err)

		if err == nil {
			loadBalancer, err := ts.driver.SelectOneLoadBalancer(testCtx, test.lbIDInput)
			ts.Equal(test.err, err)
			ts.Equal(test.lbIDInput, loadBalancer.LbID)
			ts.Equal(test.expectedUsersJSON, loadBalancer.Users)
			ts.NotEmpty(loadBalancer.CreatedAt)
			ts.NotEmpty(loadBalancer.UpdatedAt)

			users := []types.UserAccess{}
			err = json.Unmarshal(loadBalancer.Users, &users)
			ts.NoError(err)
			ts.Equal(test.expectedUsers, users)
		}
	}
}

func (ts *PGDriverTestSuite) Test_UpdateLoadBalancer() {
	tests := []struct {
		name                string
		loadBalancerID      string
		loadBalancerUpdate  *types.UpdateLoadBalancer
		expectedAfterUpdate SelectOneLoadBalancerRow
		err                 error
	}{
		{
			name:           "Should update a single load balancer successfully with all fields",
			loadBalancerID: "test_lb_34987u329rfn23f",
			loadBalancerUpdate: &types.UpdateLoadBalancer{
				Name: "vipr_app_updated",
				StickyOptions: &types.UpdateStickyOptions{
					Duration:      "100",
					StickyOrigins: []string{"chrome-extension://", "test-ext://"},
					StickyMax:     500,
					Stickiness:    boolPointer(false),
				},
			},
			expectedAfterUpdate: SelectOneLoadBalancerRow{
				Name:       sql.NullString{Valid: true, String: "vipr_app_updated"},
				Duration:   sql.NullString{Valid: true, String: "100"},
				StickyMax:  sql.NullInt32{Valid: true, Int32: 500},
				Stickiness: sql.NullBool{Valid: true, Bool: false},
				Origins:    []string{"chrome-extension://", "test-ext://"},
			},
			err: nil,
		},
		{
			name:           "Should update a single load balancer successfully with only some sticky options fields",
			loadBalancerID: "test_lb_3890ru23jfi32fj",
			loadBalancerUpdate: &types.UpdateLoadBalancer{
				Name: "vipr_app_updated_2",
				StickyOptions: &types.UpdateStickyOptions{
					Duration: "100",
				},
			},
			expectedAfterUpdate: SelectOneLoadBalancerRow{
				Name:       sql.NullString{Valid: true, String: "vipr_app_updated_2"},
				Duration:   sql.NullString{Valid: true, String: "100"},
				StickyMax:  sql.NullInt32{Valid: true, Int32: 400},
				Stickiness: sql.NullBool{Valid: true, Bool: true},
				Origins:    []string{"chrome-extension://"},
			},
			err: nil,
		},
		{
			name:           "Should update a single load balancer successfully with no sticky options fields",
			loadBalancerID: "test_lb_34gg4g43g34g5hh",
			loadBalancerUpdate: &types.UpdateLoadBalancer{
				Name: "vipr_app_updated_3",
			},
			expectedAfterUpdate: SelectOneLoadBalancerRow{
				Name:       sql.NullString{Valid: true, String: "vipr_app_updated_3"},
				Duration:   sql.NullString{Valid: true, String: "20"},
				StickyMax:  sql.NullInt32{Valid: true, Int32: 600},
				Stickiness: sql.NullBool{Valid: true, Bool: false},
				Origins:    []string{"test-extension://", "test-extension2://"},
			},
			err: nil,
		},
		{
			name:           "Should update a single load balancer successfully with only sticky options origin field",
			loadBalancerID: "test_lb_34gg4g43g34g5hh",
			loadBalancerUpdate: &types.UpdateLoadBalancer{
				StickyOptions: &types.UpdateStickyOptions{
					StickyOrigins: []string{"chrome-extension://", "test-ext://"},
				},
			},
			expectedAfterUpdate: SelectOneLoadBalancerRow{
				Name:       sql.NullString{Valid: true, String: "vipr_app_updated_3"},
				Duration:   sql.NullString{Valid: true, String: "20"},
				StickyMax:  sql.NullInt32{Valid: true, Int32: 600},
				Stickiness: sql.NullBool{Valid: true, Bool: false},
				Origins:    []string{"chrome-extension://", "test-ext://"},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		_, err := ts.driver.SelectOneLoadBalancer(testCtx, test.loadBalancerID)
		ts.Equal(test.err, err)

		err = ts.driver.UpdateLoadBalancer(testCtx, test.loadBalancerID, test.loadBalancerUpdate)
		ts.Equal(test.err, err)

		lbAfterUpdate, err := ts.driver.SelectOneLoadBalancer(testCtx, test.loadBalancerID)
		ts.Equal(test.err, err)
		ts.Equal(test.expectedAfterUpdate.Name, lbAfterUpdate.Name)
		ts.Equal(test.expectedAfterUpdate.Duration, lbAfterUpdate.Duration)
		ts.Equal(test.expectedAfterUpdate.Origins, lbAfterUpdate.Origins)
		ts.Equal(test.expectedAfterUpdate.StickyMax, lbAfterUpdate.StickyMax)
		ts.Equal(test.expectedAfterUpdate.Stickiness, lbAfterUpdate.Stickiness)
	}
}

func (ts *PGDriverTestSuite) Test_UpdateUserAccessRole() {
	tests := []struct {
		name                   string
		lbIDInput, userIDInput string
		userRoleInput          types.RoleName
		expectedUsersJSON      json.RawMessage
		expectedUsers          []types.UserAccess
		err                    error
	}{
		{
			name:              "Should update the RoleName of a UserAccess row for a LoadBalancer with correct input",
			lbIDInput:         "test_lb_3890ru23jfi32fj",
			userIDInput:       "test_user_admin5678",
			userRoleInput:     types.RoleMember,
			expectedUsersJSON: json.RawMessage(`[{"email": "owner2@test.com", "userID": "test_user_04228205bd261a", "accepted": true, "roleName": "OWNER"}, {"email": "admin2@test.com", "userID": "test_user_admin5678", "accepted": true, "roleName": "MEMBER"}]`),
			expectedUsers: []types.UserAccess{
				{
					UserID:   "test_user_04228205bd261a",
					RoleName: types.RoleOwner,
					Email:    "owner2@test.com",
					Accepted: true,
				},
				{
					UserID:   "test_user_admin5678",
					RoleName: types.RoleMember,
					Email:    "admin2@test.com",
					Accepted: true,
				},
			},
			err: nil,
		},
		{
			name:              "Should update the RoleName of a UserAccess row back to the original value for a LoadBalancer with correct input",
			lbIDInput:         "test_lb_3890ru23jfi32fj",
			userIDInput:       "test_user_admin5678",
			userRoleInput:     types.RoleAdmin,
			expectedUsersJSON: json.RawMessage(`[{"email": "owner2@test.com", "userID": "test_user_04228205bd261a", "accepted": true, "roleName": "OWNER"}, {"email": "admin2@test.com", "userID": "test_user_admin5678", "accepted": true, "roleName": "ADMIN"}]`),
			expectedUsers: []types.UserAccess{
				{
					UserID:   "test_user_04228205bd261a",
					RoleName: types.RoleOwner,
					Email:    "owner2@test.com",
					Accepted: true,
				},
				{
					UserID:   "test_user_admin5678",
					RoleName: types.RoleAdmin,
					Email:    "admin2@test.com",
					Accepted: true,
				},
			},
			err: nil,
		},
		{
			name:          "Should fail if attempting to update User to owner role",
			lbIDInput:     "test_lb_3890ru23jfi32fj",
			userIDInput:   "test_user_admin5678",
			userRoleInput: types.RoleOwner,
			err:           ErrCannotSetToOwner,
		},
		{
			name:        "Should fail if user ID not provided",
			lbIDInput:   "test_lb_34gg4g43g34g5hh",
			userIDInput: "",
			err:         ErrMissingID,
		},
		{
			name:        "Should fail if lb ID not provided",
			lbIDInput:   "",
			userIDInput: "test_user_member5678",
			err:         ErrMissingID,
		},
	}

	for _, test := range tests {
		err := ts.driver.UpdateUserAccessRole(testCtx, test.userIDInput, test.lbIDInput, test.userRoleInput)
		ts.Equal(test.err, err)

		if err == nil {
			loadBalancer, err := ts.driver.SelectOneLoadBalancer(testCtx, test.lbIDInput)
			ts.Equal(test.err, err)
			ts.Equal(test.lbIDInput, loadBalancer.LbID)
			ts.Equal(test.expectedUsersJSON, loadBalancer.Users)
			ts.NotEmpty(loadBalancer.CreatedAt)
			ts.NotEmpty(loadBalancer.UpdatedAt)

			users := []types.UserAccess{}
			err = json.Unmarshal(loadBalancer.Users, &users)
			ts.NoError(err)
			ts.Equal(test.expectedUsers, users)
		}
	}
}

func (ts *PGDriverTestSuite) Test_RemoveLoadBalancer() {
	tests := []struct {
		name           string
		loadBalancerID string
		err            error
	}{
		{
			name:           "Should remove a single load balancer successfully with correct input",
			loadBalancerID: "test_lb_34gg4g43g34g5hh",
			err:            nil,
		},
	}

	for _, test := range tests {
		err := ts.driver.RemoveLoadBalancer(testCtx, test.loadBalancerID)
		ts.Equal(test.err, err)

		lbAfterRemove, err := ts.driver.SelectOneLoadBalancer(testCtx, test.loadBalancerID)
		ts.Equal(test.err, err)
		ts.Empty(lbAfterRemove.UserID.String)
	}
}

func (ts *PGDriverTestSuite) Test_RemoveUserAccess() {
	tests := []struct {
		name                                     string
		lbIDInput, userIDInput                   string
		usersBeforeDeleteJSON, expectedUsersJSON json.RawMessage
		expectedUsers                            []types.UserAccess
		err                                      error
	}{
		{
			name:                  "Should delete a UserAccess row for a LoadBalancer with correct input",
			lbIDInput:             "test_lb_34gg4g43g34g5hh",
			userIDInput:           "test_user_member5678",
			usersBeforeDeleteJSON: json.RawMessage(`[{"email": "owner3@test.com", "userID": "test_user_redirect233344", "accepted": true, "roleName": "OWNER"}, {"email": "member2@test.com", "userID": "test_user_member5678", "accepted": true, "roleName": "MEMBER"}]`),
			expectedUsersJSON:     json.RawMessage(`[{"email": "owner3@test.com", "userID": "test_user_redirect233344", "accepted": true, "roleName": "OWNER"}]`),
			expectedUsers: []types.UserAccess{
				{
					UserID:   "test_user_redirect233344",
					RoleName: types.RoleOwner,
					Email:    "owner3@test.com",
					Accepted: true,
				},
			},
			err: nil,
		},
		{
			name:        "Should fail if user ID not provided",
			lbIDInput:   "test_lb_34gg4g43g34g5hh",
			userIDInput: "",
			err:         ErrMissingID,
		},
		{
			name:        "Should fail if lb ID not provided",
			lbIDInput:   "",
			userIDInput: "test_user_member5678",
			err:         ErrMissingID,
		},
	}

	for _, test := range tests {
		if test.err == nil {
			loadBalancerBefore, err := ts.driver.SelectOneLoadBalancer(testCtx, test.lbIDInput)
			ts.NoError(err)
			ts.Equal(test.usersBeforeDeleteJSON, loadBalancerBefore.Users)
		}

		err := ts.driver.RemoveUserAccess(testCtx, test.userIDInput, test.lbIDInput)
		ts.Equal(test.err, err)

		if test.err == nil {
			loadBalancer, err := ts.driver.SelectOneLoadBalancer(testCtx, test.lbIDInput)
			ts.Equal(test.err, err)
			ts.Equal(test.lbIDInput, loadBalancer.LbID)
			ts.Equal(test.expectedUsersJSON, loadBalancer.Users)
			ts.NotEmpty(loadBalancer.CreatedAt)
			ts.NotEmpty(loadBalancer.UpdatedAt)

			users := []types.UserAccess{}
			err = json.Unmarshal(loadBalancer.Users, &users)
			ts.NoError(err)
			ts.Equal(test.expectedUsers, users)
		}
	}
}
