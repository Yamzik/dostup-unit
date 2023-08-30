package main

import (
	"core"
	"interop"
	"net/http"
	"strconv"
	"util"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	peerRepo        *core.PeerRepository
	peerService     *core.PeerService
	mementoProvider *MementoProvider
}

func (uc *UserController) New(peerService *core.PeerService,
	memProvider *MementoProvider) *UserController {

	uc.peerService = peerService
	uc.mementoProvider = memProvider
	return uc
}
func (uc *UserController) Use(r *gin.Engine) *UserController {
	r.GET("/user/all", uc.getAllPeers)
	r.GET("/user", uc.getPeersForUser)
	r.GET("/user/specific", uc.getPeerForUser)
	r.GET("/user/download", uc.downloadConfig)
	r.POST("/user", uc.addPeerForUser)
	r.PATCH("/user/disable", uc.disablePeer)
	r.PATCH("/user/enable", uc.enablePeer)
	r.DELETE("/user", uc.deletePeerForUser)
	return uc
}

// getAllPeers
// @Summary Returns all peers registered on unit
// @Produce json
// @Param	auth   header    string  false  "auth password"
// @Success 200 {array} core.PeerSM
// @Router /api/user/all [get]
func (uc *UserController) getAllPeers(c *gin.Context) {
	peers := uc.peerService.GetAllPeers()
	c.JSON(http.StatusOK, peers)
}

// getPeersForUser
// @Summary Returns all peers on unit that belong to user
// @Param	tid    query     uint64  false  "user telegram id"
// @Param	auth   header    string  false  "auth password"
// @Produce json
// @Success 200 {array} core.PeerSM
// @Failure 400
// @Router /api/user [get]
func (uc *UserController) getPeersForUser(c *gin.Context) {
	tid_s := c.Query("tid")
	tid, err := strconv.ParseUint(tid_s, 10, 64)
	if err != nil {
		err = util.DErr(util.InvalidParameter, err.Error()).
			SetMessage("Expected tid parameter to be unsigned int of 64 bits")
		c.Error(err)
		return
	}
	peers := uc.peerService.GetPeers(tid)
	c.JSON(http.StatusOK, peers)
}

// getPeerForUser
// @Summary Returns specific peer on unit that belong to user
// @Param	tid    query     uint64  false  "user telegram id"
// @Param	auth   header    string  false  "auth password"
// @Produce json
// @Success 200 {object} core.PeerSM
// @Failure 409
// @Failure 400
// @Router /api/user/specific [get]
func (uc *UserController) getPeerForUser(c *gin.Context) {
	tid_s := c.Query("tid")
	pub := c.Query("pub")
	tid, err := strconv.ParseUint(tid_s, 10, 64)
	if err != nil {
		err = util.DErr(util.InvalidParameter, err.Error()).
			SetMessage("Expected tid parameter to be unsigned int of 64 bits")
		c.Error(err)
		return
	}
	peer, err := uc.peerService.GetPeer(tid, pub)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, peer)
}

// addPeerForUser
// @Summary Adds a peer to user by telegram id and returns created peer
// @Param	tid    query     uint64  false  "user telegram id"
// @Produce json
// @Param	auth   header    string  false  "auth password"
// @Success 200 {object} core.PeerSM
// @Failure 400
// @Failure 404
// @Failure 503
// @Router /api/user [post]
func (uc *UserController) addPeerForUser(c *gin.Context) {
	tid_s := c.Query("tid")
	tid, err := strconv.ParseUint(tid_s, 10, 64)
	if err != nil {
		err = util.DErr(util.InvalidParameter, err.Error()).
			SetMessage("Expected tid parameter to be unsigned int of 64 bits")
		c.Error(err)
		return
	}
	// todo add transaction to prevent adding but not saving
	peer, err := uc.peerService.Add(tid)
	if err != nil {
		c.Error(err)
		return
	}
	err = uc.sync()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, peer)
}

// deletePeerForUser
// @Summary Deletes peer of user by telegram id and peer public key
// @Param	tid    query     uint64  false  "user telegram id"
// @Param	pub    query     string  false  "peer public key"
// @Param	auth   header    string  false  "auth password"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 503
// @Router /api/user [delete]
func (uc *UserController) deletePeerForUser(c *gin.Context) {
	tid_s := c.Query("tid")
	pub := c.Query("pub")
	tid, err := strconv.ParseUint(tid_s, 10, 64)
	if err != nil {
		err = util.DErr(util.InvalidParameter, err.Error()).
			SetMessage("Expected tid parameter to be unsigned int of 64 bits")
		c.Error(err)
		return
	}
	// todo add transaction to prevent adding but not saving
	err = uc.peerService.Delete(tid, pub)
	if err != nil {
		c.Error(err)
		return
	}
	err = uc.sync()
	if err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusOK)
}

// downloadConfig
// @Summary Builds ready-to-use client as string
// @Param	tid    query     uint64  false  "user telegram id"
// @Param	pub    query     string  false  "peer public key"
// @Param	auth   header    string  false  "auth password"
// @Success 200 {string} string
// @Failure 400
// @Failure 404
// @Router /api/user/download [get]
func (uc *UserController) downloadConfig(c *gin.Context) {
	tid_s := c.Query("tid")
	pub := c.Query("pub")
	tid, err := strconv.ParseUint(tid_s, 10, 64)
	if err != nil {
		err = util.DErr(util.InvalidParameter, err.Error()).
			SetMessage("Expected tid parameter to be unsigned int of 64 bits")
		c.Error(err)
		return
	}

	client, err := uc.peerService.GetPeerInterface(tid, pub)
	if err != nil {
		c.Error(err)
		return
	}

	c.String(http.StatusOK, client)
}

// enablePeer
// @Summary Enables peer by telegram id and peer public key
// @Param	tid    query     uint64  false  "user telegram id"
// @Param	pub    query     string  false  "peer public key"
// @Param	auth   header    string  false  "auth password"
// @Success 200 {object} core.PeerSM
// @Failure 400
// @Failure 404
// @Failure 500
// @Failure 503
// @Router /api/user/enable [patch]
func (uc *UserController) enablePeer(c *gin.Context) {
	tid_s := c.Query("tid")
	pub := c.Query("pub")
	tid, err := strconv.ParseUint(tid_s, 10, 64)
	if err != nil {
		err = util.DErr(util.InvalidParameter, err.Error()).
			SetMessage("Expected tid parameter to be unsigned int of 64 bits")
		c.Error(err)
		return
	}

	peer, err := uc.peerService.EnablePeer(tid, pub)
	if err != nil {
		c.Error(err)
		return
	}
	err = uc.sync()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, peer)
}

// disablePeer
// @Summary Disables peer by telegram id and peer public key
// @Param	tid    query     uint64  false  "user telegram id"
// @Param	pub    query     string  false  "peer public key"
// @Param	auth   header    string  false  "auth password"
// @Success 200 {object} core.PeerSM
// @Failure 400
// @Failure 404
// @Failure 500
// @Failure 503
// @Router /api/user/disable [patch]
func (uc *UserController) disablePeer(c *gin.Context) {
	tid_s := c.Query("tid")
	pub := c.Query("pub")
	tid, err := strconv.ParseUint(tid_s, 10, 64)
	if err != nil {
		err = util.DErr(util.InvalidParameter, err.Error()).
			SetMessage("Expected tid parameter to be unsigned int of 64 bits")
		c.Error(err)
		return
	}

	peer, err := uc.peerService.DisablePeer(tid, pub)
	if err != nil {
		c.Error(err)
		return
	}
	err = uc.sync()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, peer)
}

func (uc *UserController) sync() error {
	mem := uc.peerService.GetMemento()
	err := uc.mementoProvider.Save(mem)
	if err != nil {
		return err
	}
	config, err := uc.peerService.MakeConfig()
	if err != nil {
		return err
	}
	err = uc.mementoProvider.SaveConfig(config)
	if err != nil {
		return err
	}
	err = interop.SyncConf()
	if err != nil {
		return err
	}
	return nil
}
