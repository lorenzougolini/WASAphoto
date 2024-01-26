package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	rt.router.POST("/session", rt.userLogin)

	rt.router.GET("/users/:username", rt.getUserProfile)
	rt.router.PUT("/users/:username", rt.setUsername)

	rt.router.POST("/users/:username/photos", rt.uploadNewPhoto)
	rt.router.DELETE("/users/:username/photos/:photoid", rt.deletePhoto)

	rt.router.PUT("/users/:username/followers/:followedUsername", rt.followUser)
	rt.router.DELETE("/users/:username/followers/:unfollowedUsername", rt.unfollowUser)

	rt.router.PUT("/users/:username/banned/:bannedUsername", rt.banUser)
	rt.router.DELETE("/users/:username/banned/:unbannedUsername", rt.unbanUser)

	rt.router.POST("/photos/:photoid/likes", rt.likePhoto)
	rt.router.DELETE("/photos/:photoid/likes/:likeid", rt.unlikePhoto)

	rt.router.POST("/photos/:photoid/comments", rt.commentPhoto)
	rt.router.DELETE("/photos/:photoid/comments/:commentid", rt.uncommentPhoto)

	rt.router.GET("/stream", rt.getStream)

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
