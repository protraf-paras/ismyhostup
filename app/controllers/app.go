package controllers

import (
	"github.com/revel/revel"
	"ismyhostup/app"
)

type App struct {
	*revel.Controller
}

type Event struct {
	Id 			int
	Canonical 	string
	URL			string
	ASN 		uint16
	Status		string
	Date		uint32
}

func (c App) Index() (res revel.Result) {
	res = c.Render()

	ev, err := app.DB.Query("SELECT events.id, hosts.canonical, hosts.url, hosts.asn, events.status, events.date_detected FROM events LEFT JOIN hosts ON events.host_id = hosts.id ORDER BY events.date_detected DESC LIMIT 16")
	if err != nil {
		revel.ERROR.Printf("%s\n", err.Error())
		return
	}

	events := make([]*Event, 0)
	for ev.Next() {
		e := &Event{}
		ev.Scan(&e.Id, &e.Canonical, &e.URL, &e.ASN, &e.Status, &e.Date)
		events = append(events, e)
	}
	
	c.RenderArgs["events"] = events
	return
}

func (c App) Host() revel.Result {
	h, err := app.DB.Query("SELECT id, canonical, asn FROM hosts WHERE url = ?", c.Params.Get("host"))
	if err != nil {
		c.Flash.Error(err.Error())
		return c.Redirect(App.Index)
	}
	if !h.Next() {
		c.Flash.Error("No such host exists! To submit a request, please email paras (at) protrafsolutions (d0t) com")
		return c.Redirect(App.Index)
	}
	var hostId int
	var canonical string
	var asn uint16
	h.Scan(&hostId, &canonical, &asn)

	ev, err := app.DB.Query("SELECT id, status, date_detected FROM events WHERE host_id = ? ORDER BY date_detected DESC", hostId)
	if err != nil {
		c.Flash.Error(err.Error())
		return c.Redirect(App.Index)
	}

	events := make([]*Event, 0)
	for ev.Next() {
		e := &Event{-1, canonical, c.Params.Get("host"), asn, "", uint32(0)}
		ev.Scan(&e.Id, &e.Status, &e.Date)
		events = append(events, e)
	}

	c.RenderArgs["events"] = events
	return c.Render()
}
