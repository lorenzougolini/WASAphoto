package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	rt.router.GET("/lists", rt.listAll)

	rt.router.POST("/session", rt.userLogin)
	rt.router.GET("/users/:userid", rt.getUserProfile)
	rt.router.PUT("/users/:userid", rt.setUsername)

	rt.router.POST("/users/:userid/photos", rt.uploadNewPhoto)
	rt.router.DELETE("/users/:userid/photos/:photoid", rt.deletePhoto)

	rt.router.PUT("/users/:userid/followers/:username", rt.followUser)
	rt.router.DELETE("/users/:userid/followers/:username", rt.unfollowUser)

	rt.router.PUT("/users/:userid/banned/:username", rt.banUser)
	rt.router.DELETE("/users/:userid/banned/:username", rt.unbanUser)

	rt.router.POST("/photos/:photoid/likes", rt.likePhoto)
	rt.router.DELETE("/photos/:photoid/likes/:likeid", rt.unlikePhoto)

	// rt.router.POST("/photos/:photoid/comments", rt.commentPhoto)
	// rt.router.DELETE("/photos/:photoid/comments/:commentid", rt.uncommentPhoto)

	// rt.router.GET("/stream", rt.getStream)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
