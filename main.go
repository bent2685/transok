package main

import (
	"context"
	"log"

	"transok/backend/services"
	"transok/backend/utils/logger"
	"transok/backend/utils/mdns"
	mdns_handlers "transok/backend/utils/mdns/handlers"
	"transok/desktop"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/plugins/updater/native"
	"github.com/egoist/mygo/ui"
)

func main() {
	logger.InitLogger()
	sysSvc := services.System()
	discoverSvc := services.GetDiscoverService()
	shareSvc := services.Share()
	storageSvc := services.Storage()

	appInfo := sysSvc.GetAppInfo()
	sysSvc.Start(appInfo["version"])
	storageSvc.Init(context.Background())
	discoverSvc.Start()
	shareSvc.Start(context.Background())

	/* Subscribe to mDNS messages */
	mdns.GetDispatcher().Subscribe(mdns_handlers.GetDiscoverHandler())
	mdns.GetDispatcher().Subscribe(mdns_handlers.NewPingHandler())

	mygo.Use(native.Plugin)

	application := desktop.New()

	mygo.App.WhenReady(func() {
		win := mygo.NewWindow(mygo.WindowOptions{
			Title:         appInfo["name"],
			Width:         380,
			Height:        620,
			MinWidth:      380,
			MaxWidth:      380,
			MinHeight:     620,
			MaxHeight:     620,
			DisableResize: true,
			TitleBarStyle: mygo.TitleBarHidden,
			Content:       ui.View(application.View),
		})
		application.Attach(win)

		application.SetDark(mygo.Theme.IsDark())
		mygo.Theme.OnUpdated(func() { application.SetDark(mygo.Theme.IsDark()) })
	})

	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
