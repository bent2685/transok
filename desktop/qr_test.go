package desktop

import "testing"

func TestQRBitmap(t *testing.T) {
	bmp, err := qrBitmap("http://192.168.50.22:9482")
	if err != nil {
		t.Fatalf("生成二维码失败：%v", err)
	}
	if bmp == nil {
		t.Fatal("二维码位图为空")
	}
}

func TestShowQRClosesShareDrawer(t *testing.T) {
	a := &App{shareOpen: true}
	a.showQR("http://192.168.50.22:9482")

	if !a.qrOpen || a.shareOpen {
		t.Fatalf("qrOpen=%v shareOpen=%v，期望二维码打开、分享抽屉收起", a.qrOpen, a.shareOpen)
	}
	if a.qrURL != "http://192.168.50.22:9482" || a.qrImage == nil {
		t.Fatalf("二维码状态未写入：url=%q image=%v", a.qrURL, a.qrImage)
	}
}
