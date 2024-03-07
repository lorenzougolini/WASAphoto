package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	rt.router.POST("/session", rt.wrap(rt.userLogin))

	rt.router.GET("/users/:username", rt.wrap(rt.getUserProfile))
	rt.router.PUT("/users/:username", rt.wrap(rt.setUsername))

	rt.router.POST("/users/:username/photos", rt.wrap(rt.uploadNewPhoto))
	rt.router.DELETE("/users/:username/photos/:photoid", rt.wrap(rt.deletePhoto))

	rt.router.PUT("/users/:username/followers/:followedUsername", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:username/followers/:unfollowedUsername", rt.wrap(rt.unfollowUser))

	rt.router.PUT("/users/:username/banned/:bannedUsername", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:username/banned/:unbannedUsername", rt.wrap(rt.unbanUser))

	rt.router.POST("/photos/:photoid/likes", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/photos/:photoid/likes/:likeid", rt.wrap(rt.unlikePhoto))

	rt.router.POST("/photos/:photoid/comments", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/photos/:photoid/comments/:commentid", rt.wrap(rt.uncommentPhoto))

	rt.router.GET("/stream", rt.wrap(rt.getStream))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
