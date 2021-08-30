package main

import (
	"context"
	"errors"
	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/entities"
	"github.com/Kong/go-pdk/server"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"strings"
)

// Version 定义版本
var Version = "0.1"

// Priority 优先级
var Priority = 777

// define the redis client
var rdb *redis.Client

// default response header
var header = map[string][]string{"Content-Type": {"application/json"}}

func main() {
	server.StartServer(New, Version, Priority)
}

type Config struct {
	Address     string
	Password    string
	DB          int
	TokenHeader string `json:"token_header"`
	MidHeader   string `json:"mid_header"`
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Access(kong *pdk.PDK) {

	tokenString, err := kong.Request.GetHeader(conf.TokenHeader)
	if err != nil {
		kong.Response.Exit(401, "token is not allowed empty ", header)
		return
	}

	router, err := kong.Router.GetRoute()
	if err != nil {
		kong.Response.Exit(401, "router does not matched ", header)
		return
	}

	userid, err := parseJwt(tokenString)
	if err != nil || len(userid) == 0 {
		kong.Response.Exit(401, err.Error(), header)
		return
	}

	if !conf.checkAuth(kong, router, userid) {
		kong.Response.Exit(403, "You have no permission", header)
	}
}

func (conf Config) checkAuth(kong *pdk.PDK, route entities.Route, userId string) bool {
	if rdb == nil {
		rdb = redis.NewClient(&redis.Options{
			Addr:     conf.Address,
			Password: conf.Password,
			DB:       conf.DB,
		})
	}

	//  config router
	var ConfigRoutes []string
	for _, path := range route.Paths {

		if len(route.Methods) > 0 {
			for _, method := range route.Methods {
				ConfigRoutes = append(ConfigRoutes, strings.ToUpper(method)+path)
			}
		} else if method, ok := kong.Request.GetMethod(); ok == nil {
			ConfigRoutes = append(ConfigRoutes, strings.ToUpper(method)+path)
		} else {
			return false
		}
	}

	midEncryption, _ := kong.Request.GetHeader(conf.MidHeader)
	if len(midEncryption) > 0 && !checkAuthByMid(midEncryption, ConfigRoutes) {
		return false
	}

	return checkAuthByUserid(userId, ConfigRoutes)
}

func checkAuthByUserid(userId string, ConfigRoutes []string) bool {

	userPermissionsCache, err := rdb.SMembers(context.Background(), userId).Result()
	if err != nil {
		return false
	}

	return NewSet(userPermissionsCache...).Contains(ConfigRoutes...)
}

func checkAuthByMid(midEncryption string, ConfigRoutes []string) bool {

	// mid, err := rdb.Get(context.Background(), midEncryption).Result()
	// if err != nil {
	// 	return false
	// }

	midPermissionsCache, err := rdb.SMembers(context.Background(), midEncryption).Result()
	if err != nil {
		return false
	}

	return NewSet(midPermissionsCache...).Contains(ConfigRoutes...)
}

/*
 *
 *  @Description: parse jwt without verify
 *  @param tokenString
 *  @return string
 *  @return error
 *
 **/
func parseJwt(tokenString string) (string, error) {

	token, _, _ := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if token == nil {
		return "", errors.New("userid is not allowed empty")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if userid, ok := claims["userid"]; ok {

			return userid.(string), nil
		}
	}

	return "", errors.New("userid is not allowed empty")
}
