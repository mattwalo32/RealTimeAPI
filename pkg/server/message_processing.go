package server

import (
	"github.com/mattwalo32/RealTimeAPI/internal/messages"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (router *MessageRouter) processMessage(msg messages.Message) {
	switch msg.GetMessageType() {
	case messages.MESSAGE_ACKNOWLEDGE:
		router.processAcknowledge(msg.(*messages.AcknowledgementMessage))
	case messages.MESSAGE_JOIN_SERVER:
		router.processJoinServer(msg.(*messages.JoinServerMessage))
	case messages.MESSAGE_FIND_ROOM:
		// TODO: Process this message
	default:
		log.WithFields(log.Fields{
			"messageType": msg.GetMessageType(),
		}).Warn("Router got unkown message type")
		return
	}

	if messages.DoesMessageRequireResponse(msg) {
		router.acknowledgeMessage(msg)
	}

	associableMsg, isRoutable := msg.(messages.RoutableMessage)
	if isRoutable {
		router.routeMessage(msg, associableMsg.GetClientID())
	}
}

func (router *MessageRouter) removeMessageTimer(msg *messages.AcknowledgementMessage) {
	router.lock.Lock()
	defer router.lock.Unlock()
	evtID, evtExists := router.messageRetryEventIDs[msg.AcknowledgedMessageID]
	if !evtExists {
		log.WithFields(log.Fields{
			"ID": msg.AcknowledgedMessageID,
		}).Debug("Acknowledged message without a timer")
		return
	}

	log.WithFields(log.Fields{
		"ID": msg.AcknowledgedMessageID,
	}).Debug("Recieved acknowledgement for message")
	router.timer.RemoveEvent(evtID)
	delete(router.messageRetryEventIDs, msg.AcknowledgedMessageID)
}

func (router *MessageRouter) acknowledgeMessage(msg messages.Message) {
	ackMessage := &messages.AcknowledgementMessage{
		SourceAddr: *router.udpManager.GetUDPAddr(),
		DestAddr: msg.GetSource(),
		MessageID: uuid.New(),
		AcknowledgedMessageID: msg.GetID(),
	}

	router.SendMessageUnreliably(ackMessage)
}

func (router *MessageRouter) routeMessage(msg messages.Message, clientID uuid.UUID) {
	routerLogger := log.WithFields(log.Fields{
		"ClientID": clientID,
		"Message": msg,
	})

	client, doesClientExist := router.clients[clientID]
	if !doesClientExist {
		routerLogger.Warn("Attempted to route message to non-existent client")
		return
	}

	roomID := client.RoomID
	room, doesRoomExist := router.rooms[roomID]
	if !doesRoomExist {
		routerLogger.WithFields(log.Fields{
			"RoomID": roomID,
		}).Warn("Attmpted to route message to non-existent room")
		return
	}

	room.config.EventHandler.OnMessageRecieved(msg)
}