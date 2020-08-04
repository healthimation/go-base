package token

// Permission can validate if a validated token has a given permission
type Permission struct {
	id         string
	role       string
	userAccess map[string]interface{}
	logger     Logger
	verified   bool
}

type grantClaim struct {
	grants string
	days   int64
}

func (p *Permission) log(userIDs []string, roles []string, result bool) {
	if p.logger != nil {
		p.logger.Printf("user (%s - %s) attempted to access data for users (%v) as a %s. Access granted (%v)", p.id, p.role, userIDs, roles, result)
	}
}

func (p *Permission) logGrantAccess(msg string) {
	if p.logger != nil {
		p.logger.Printf("%s", msg)
	}
}

// GetUserID will return the user's id
func (p *Permission) GetUserID() string {
	return p.id
}

// IsUser returns true if this permission is for the provided user
func (p *Permission) IsUser(userID string) bool {
	return p.id == userID
}

// IsRole returns true if this permission is for the provided role
func (p *Permission) IsRole(role string) bool {
	return p.role == role
}

// IsVerified returns true if the user has verified their email address
func (p *Permission) IsVerified() bool {
	return p.verified
}

// CanAccessAs returns true if this permission is the role require and has access to the users proviced
func (p *Permission) CanAccessAs(userIDs []string, roles []string) (result bool) {
	result = false
	defer func() {
		p.log(userIDs, roles, result)
	}()

	for _, role := range roles {
		if p.role == role {
			switch role {
			case RoleAdmin:
				fallthrough
			case RoleCS:
				// CS and admin have access to all users
				result = true
			case RoleRN:
				fallthrough
			case RoleCoach:
				// RN and coach only have access to specific users
				result = true
				for _, u := range userIDs {
					if _, ok := p.userAccess[u]; !ok {
						// they don't have access to this user
						result = false
						break
					}
				}
			}

			// if we found a role they have access to return true.
			if result {
				return result
			}
			//else continue to loop
		}
	}
	return false
}

// // GrantAccess takes a grant string and verifies it is present in the claims and
// // that the user is not outside of the granted time.
// func (p *Permission) GrantAccess(grantToVerify string) bool {

// 	var gc *grantClaim
// 	for _, grant := range p.grants {
// 		grantMap, gmOK := grant.(map[string]interface{})
// 		if !gmOK {
// 			p.logGrantAccess(fmt.Sprintf("Grant is not the proper type: %#v", grant))
// 			return false
// 		}

// 		g, gOK := grantMap["g"]
// 		d, dOK := grantMap["d"]

// 		if !gOK || !dOK {
// 			p.logGrantAccess("Grant string or grant day not set in grant")
// 			return false
// 		}

// 		grantStr, gOK := g.(string)
// 		if !gOK {
// 			p.logGrantAccess("Grant string is not the proper type")
// 			return false
// 		}

// 		if grantStr == grantToVerify {
// 			if d == nil {
// 				//infinite grant
// 				return true
// 			}

// 			daysFloat, dOK := d.(float64)
// 			if !dOK {
// 				p.logGrantAccess("Grant days is not the proper type")
// 				return false
// 			}

// 			daysInt := int64(daysFloat)
// 			gc = &grantClaim{
// 				grants: grantStr,
// 				days:   daysInt,
// 			}
// 		}
// 	}

// 	if gc == nil {
// 		p.logGrantAccess("Grant string not found")
// 		return false
// 	}

// 	programDay, err := p.getProgramDay()
// 	if err != nil {
// 		p.logGrantAccess(fmt.Sprintf("Unable to determine program day: %v", err))
// 		return false
// 	}

// 	if programDay <= gc.days {
// 		return true
// 	}

// 	// Default to false
// 	p.logGrantAccess("Grant access fell through")
// 	return false
// }

// // StaticDataID returns the static data ID from the permission object
// func (p *Permission) StaticDataID() string {
// 	return p.staticDataID
// }

// // ProgramDay calculates the program day from the logical start date
// func (p *Permission) ProgramDay() (int64, error) {
// 	return p.getProgramDay()
// }

// // program day can be negative, indicating a start date in the future.
// func (p *Permission) getProgramDay() (int64, error) {
// 	if p.logicalStartDate == nil {
// 		return 0, fmt.Errorf("No logical start date")
// 	}

// 	lsd, err := time.Parse(time.RFC3339, *p.logicalStartDate)
// 	if err != nil {
// 		return 0, fmt.Errorf("Unable to parse time from logical start date: %v", err)
// 	}

// 	programDur := time.Now().UTC().Sub(lsd)
// 	return int64(programDur.Hours() / 24), nil
// }
