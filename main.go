package main

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"image/color"
	"io"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MythosV5Theme struct{}

func (m MythosV5Theme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 2, G: 4, B: 2, A: 255}
	case theme.ColorNameForeground:
		return color.RGBA{R: 0, G: 240, B: 40, A: 255}
	case theme.ColorNamePrimary:
		return color.RGBA{R: 0, G: 180, B: 30, A: 255}
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 10, G: 15, B: 10, A: 255}
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (m MythosV5Theme) Font(style fyne.TextStyle) fyne.Resource { return theme.DefaultTheme().Font(style) }
func (m MythosV5Theme) Size(name fyne.ThemeSizeName) float32   { return theme.DefaultTheme().Size(name) }
func (m MythosV5Theme) Icon(name fyne.ThemeIconName) fyne.Resource { return theme.DefaultTheme().Icon(name) }

type Specs struct {
	Cores int
	Tier  string
	Model string
}

func checkSpecs() Specs {
	cores := runtime.NumCPU()
	s := Specs{Cores: cores}
	if cores > THE MYTHOS BOOT SEQUENCE COMPLETED.\n>>> CONFIGURING CLIENT OS V5.1 INTERFACE...\n")
	logTerminal.Wrapping = fyne.TextWrapWord

	pushLog := func(info string) {
		ts := time.Now().Format("15:04:05")
		logTerminal.SetText(logTerminal.Text + fmt.Sprintf("[%s] %s\n", ts, info))
	}

	sessionHistories := make(map[string][]OllamaMessage)
	currentSessionID := "SESSION_ROOT"
	sessionHistories[currentSessionID] = []OllamaMessage{}

	chatOutput := widget.NewMultiLineEntry()
	chatOutput.SetText(">>> SYSTEM STANDBY. AWAITING AI FRAME CAPTURE...\n")
	chatInput := widget.NewEntry()
	chatInput.SetPlaceHolder("Inject user query into sandbox channel...")

	refreshChatUI := func() {
		var buf bytes.Buffer
		for _, msg := range sessionHistories[currentSessionID] {
			if msg.Role == "user" {
				buf.WriteString("INJECTED_USER: " + msg.Content + "\n")
			} else if msg.Role == "assistant" {
				buf.WriteString("MYTHOS_CORE_RESP: " + msg.Content + "\n")
			} else if msg.Role == "system" {
				buf.WriteString("SYSTEM_CONTEXT_INJECTED: " + msg.Content + "\n")
			}
		}
		if buf.Len() == 0 {
			chatOutput.SetText(">>> SYSTEM STANDBY. AWAITING AI FRAME CAPTURE...\n")
		} else {
			chatOutput.SetText(buf.String())
		}
	}

	firewallLabel := widget.NewLabel("⚠ [ FIREWALL INTERCEPT ALERT ]\n攔截到外流網絡連線封包意圖...")
	var fwStatus string = "BLOCKED"
	fwStatusLabel := widget.NewLabel("CURRENT PROTOCOL: " + fwStatus)

	allowBtn := widget.NewButton("ALLOW 出關", func() {
		fwStatus = "ALLOWED"
		fwStatusLabel.SetText("CURRENT PROTOCOL: " + fwStatus)
		pushLog("CRITICAL: Outbound firewall dropped. Sandbox bypass authorized by user.")
	})
	blockBtn := widget.NewButton("BLOCK 攔截", func() {
		fwStatus = "BLOCKED"
		fwStatusLabel.SetText("CURRENT PROTOCOL: " + fwStatus)
		pushLog("FIREWALL SHIELD REINFORCED. Outbound attempt neutralized.")
	})

	openParamBtn := widget.NewButton("CONFIG ALL HYPERPARAMETERS", func() {
		paramWin := myApp.NewWindow("THE MYTHOS ENGINE - ADVANCED HYPERPARAMETERS")
		paramWin.Resize(fyne.NewSize(550, 550))

		sysPromptEntry := widget.NewMultiLineEntry()
		sysPromptEntry.SetText(aiCfg.SysPrompt)
		tempEntry := widget.NewEntry()
		tempEntry.SetText(fmt.Sprintf("%.2f", aiCfg.Temperature))
		topPEntry := widget.NewEntry()
		topPEntry.SetText(fmt.Sprintf("%.2f", aiCfg.TopP))
		ctxLenEntry := widget.NewEntry()
		ctxLenEntry.SetText(strconv.Itoa(aiCfg.ContextLen))
		topKEntry := widget.NewEntry()
		topKEntry.SetText(strconv.Itoa(aiCfg.TopK))
		repPenEntry := widget.NewEntry()
		repPenEntry.SetText(fmt.Sprintf("%.2f", aiCfg.RepeatPen))
		seedEntry := widget.NewEntry()
		seedEntry.SetText(strconv.Itoa(aiCfg.Seed))
		maxTokEntry := widget.NewEntry()
		maxTokEntry.SetText(strconv.Itoa(aiCfg.MaxTokens))
		ollamaURLEntry := widget.NewEntry()
		ollamaURLEntry.SetText(aiCfg.OllamaAPI)
		cryptKeyEntry := widget.NewEntry()
		cryptKeyEntry.SetText(aiCfg.SecretKey)

		form := widget.NewForm(
			widget.NewFormItem("SYSTEM_PROMPT", sysPromptEntry),
			widget.NewFormItem("TEMPERATURE", tempEntry),
			widget.NewFormItem("TOP_P", topPEntry),
			widget.NewFormItem("CONTEXT_LENGTH", ctxLenEntry),
			widget.NewFormItem("TOP_K", topKEntry),
			widget.NewFormItem("REPEAT_PENALTY", repPenEntry),
			widget.NewFormItem("SEED", seedEntry),
			widget.NewFormItem("MAX_TOKENS", maxTokEntry),
			widget.NewFormItem("OLLAMA_ENDPOINT", ollamaURLEntry),
			widget.NewFormItem("AES_KEY_32BYTE", cryptKeyEntry),
		)

		saveBtn := widget.NewButton("SAVE & INJECT", func() {
			aiCfg.SysPrompt = sysPromptEntry.Text
			aiCfg.Temperature, _ = strconv.ParseFloat(tempEntry.Text, 64)
			aiCfg.TopP, _ = strconv.ParseFloat(topPEntry.Text, 64)
			aiCfg.ContextLen, _ = strconv.Atoi(ctxLenEntry.Text)
			aiCfg.TopK, _ = strconv.Atoi(topKEntry.Text)
			aiCfg.RepeatPen, _ = strconv.ParseFloat(repPenEntry.Text, 64)
			aiCfg.Seed, _ = strconv.Atoi(seedEntry.Text)
			aiCfg.MaxTokens, _ = strconv.Atoi(maxTokEntry.Text)
			aiCfg.OllamaAPI = ollamaURLEntry.Text
			aiCfg.SecretKey = cryptKeyEntry.Text
			pushLog("ENGINE PARAMETERS MODIFIED: Real local Ollama payload model re-bound.")
			paramWin.Close()
		})

		paramWin.SetContent(container.NewBorder(
			widget.NewLabel("=== ENGINE MATRIX PARAMETERS OVERRIDE ==="),
			saveBtn, nil, nil, container.NewVScroll(form),
		))
		paramWin.Show()
	})

	makeTorDialer := func() func(ctx context.Context, network, addr string) (net.Conn, error) {
		return func(ctx context.Context, network, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout("tcp", aiCfg.TorProxyAddr, 4*time.Second)
			if err != nil {
				return nil, err
			}
			if _, err := conn.Write([]byte{0x05, 0x01, 0x00}); err != nil {
				conn.Close()
				return nil, err
			}
			buf := make([]byte, 2)
			if _, err := io.ReadFull(conn, buf); err != nil {
				conn.Close()
				return nil, fmt.Errorf("socks5 handshake failed")
			}
			host, portStr, _ := net.SplitHostPort(addr)
			port, _ := strconv.Atoi(portStr)
			reqBytes := []byte{0x05, 0x01, 0x00, 0x03, byte(len(host))}
			reqBytes = append(reqBytes, []byte(host)...)
			reqBytes = append(reqBytes, byte(port>>8), byte(port))
			if _, err := conn.Write(reqBytes); err != nil {
				conn.Close()
				return nil, err
			}
			reply := make([]byte, 4)
			if _, err := io.ReadFull(conn, reply); err != nil {
				conn.Close()
				return nil, fmt.Errorf("circuit tunnel refused")
			}
			dummy := make([]byte, 6)
			io.ReadFull(conn, dummy)
			return conn, nil
		}
	}

	modelCheck := widget.NewCheck(sysSpecs.Model+" ("+sysSpecs.Tier+")", func(b bool) {})
	modelCheck.SetChecked(true)

	darkWebTargetEntry := widget.NewEntry()
	darkWebTargetEntry.SetText("http://torproject.org")

	executeModelBtn := widget.NewButton("EXECUTE / DEPLOY", func() {
		pushLog("Analyzing Hardware... Target active tag: " + sysSpecs.Model)
		pushLog("Verifying Local Ollama Core connection on endpoint: " + aiCfg.OllamaAPI)

		go func() {
			testClient := &http.Client{Timeout: 2 * time.Second}
			resp, err := testClient.Get("http://localhost:11434/api/tags")
			if err != nil {
				pushLog("WARN: Local Ollama daemon unreachable via HTTP. Confirm port binding.")
			} else {
				pushLog("SUCCESS: Local Ollama pipeline handshake complete. Engine ready.")
				resp.Body.Close()
			}

			pushLog(">>> AUTOMATED TOR RADAR: ACTIVATED BACKGROUND CIRCUIT SCANNING...")
			for {
targetSite := darkWebTargetEntry.Text
if targetSite == "" {
targetSite = "torproject.org"
}
transport := &http.Transport{DialContext: makeTorDialer()}
client := &http.Client{Transport: transport, Timeout: 8 * time.Second}
resp, err := client.Get(targetSite)
if err != nil {
pushLog("TOR REAL CRITICAL: Physical proxy tunnel unreachable at " + aiCfg.TorProxyAddr)
} else {
body, _ := io.ReadAll(io.LimitReader(resp.Body, 128))
textCaptured := string(bytes.ReplaceAll(body, []byte("\n"), []byte(" ")))
if len(textCaptured) > 60 {
textCaptured = textCaptured[:60]
}
pushLog("[REAL TOR INGESTION SUCCESS] Packet captured.")
sessionHistories[currentSessionID] = append(sessionHistories[currentSessionID], OllamaMessage{
Role:    "system",
Content: "Real-time network intelligence update: " + textCaptured,
})
fyne.DoQueue(refreshChatUI)
resp.Body.Close()
}
time.Sleep(5 * time.Second)
}
}()
})
surfaceWebQuery := widget.NewEntry()
surfaceWebQuery.SetPlaceHolder("Enter clearnet domain URL...")
launchSurfaceBtn := widget.NewButton("LAUNCH SURFACE SCAN", func() {
target := surfaceWebQuery.Text
if target == "" {
return
}
pushLog("Executing HTTP standard crawler stack against: " + target)
go func() {
if fwStatus == "BLOCKED" {
pushLog("SCAN ERROR: Outbound network blocked by Firewall Shield.")
return
}
client := &http.Client{Timeout: 5 * time.Second}
resp, err := client.Get(target)
if err != nil {
pushLog("Network connection timeout: " + err.Error())
return
}
defer resp.Body.Close()
body, _ := io.ReadAll(io.LimitReader(resp.Body, 128))
textCaptured := string(bytes.ReplaceAll(body, []byte("\n"), []byte(" ")))
if len(textCaptured) > 60 {
textCaptured = textCaptured[:60]
}
pushLog("Surface payload captured and injected into memory block.")
sessionHistories[currentSessionID] = append(sessionHistories[currentSessionID], OllamaMessage{
Role:    "system",
Content: "Clearnet vector scan result injected: " + textCaptured,
})
fyne.DoQueue(refreshChatUI)
}()
})
sessionSelect := widget.NewSelect([]string{"SESSION_ROOT"}, func(value string) {
if value != "" {
currentSessionID = value
refreshChatUI()
pushLog("Switched execution memory frame to: " + currentSessionID)
}
})
sessionSelect.SetSelected("SESSION_ROOT")
newSessionBtn := widget.NewButton("NEW CONTEXT", func() {
newID := fmt.Sprintf("SESSION_%04D", time.Now().UnixNano()%10000)
sessionHistories[newID] = []OllamaMessage{}
var currentOptions []string
currentOptions = append(sessionSelect.Options, newID)
sessionSelect.Options = currentOptions
sessionSelect.SetSelected(newID)
pushLog("Allocated new sandboxed multi-context node: " + newID)
})
purgeSessionBtn := widget.NewButton("PURGE CONTEXT", func() {
sessionHistories[currentSessionID] = []OllamaMessage{}
refreshChatUI()
pushLog("CRITICAL: Context stack cleared for node " + currentSessionID)
})
exportSecureLogBtn := widget.NewButton("EXPORT SECURE BINARY LOG", func() {
pushLog("Initiating cryptographic storage cycle for: " + currentSessionID)
go func() {
if len(sessionHistories[currentSessionID]) == 0 {
pushLog("EXPORT REJECTED: Targeted session frame buffer is completely empty.")
return
}
rawJSON, err := json.Marshal(sessionHistories[currentSessionID])
if err != nil {
pushLog("CRYPTO ERROR: Failed to marshal target structural array.")
return
}
keyBytes := []byte(aiCfg.SecretKey)
if len(keyBytes) != 32 {
pushLog("CRYPTO ABORT: AES key vector constraints violated. Must be exactly 32 bytes.")
return
}
encryptedData, err := encryptGCM(rawJSON, keyBytes)
if err != nil {
pushLog("CRYPTO ERROR: Encryption engine pipeline fault -> " + err.Error())
return
}
fileName := fmt.Sprintf("mythos_%s_%d.enc.bin", currentSessionID, time.Now().Unix())
err = os.WriteFile(fileName, encryptedData, 0600)
if err != nil {
pushLog("DISK WRITE FAILURE: Permission block or path exception -> " + err.Error())
return
}
pushLog("SUCCESS: Secure binary frame written locally to -> " + fileName)
}()
})
sendPromptBtn := widget.NewButton("INJECT PROMPT", func() {
text := chatInput.Text
if text == "" {
return
}
sessionHistories[currentSessionID] = append(sessionHistories[currentSessionID], OllamaMessage{
Role:    "user",
Content: text,
})
refreshChatUI()
pushLog("Token pipeline processing: Transmitting payload vector to Ollama API cluster...")
chatInput.SetText("")
go func() {
fullMessages := []OllamaMessage{
{Role: "system", Content: aiCfg.SysPrompt},
}
fullMessages = append(fullMessages, sessionHistories[currentSessionID]...)
reqPayload := OllamaChatRequest{
Model:    sysSpecs.Model,
Messages: fullMessages,
Stream:   false,
Options: OllamaOptions{
Temperature:   aiCfg.Temperature,
TopP:          aiCfg.TopP,
TopK:          aiCfg.TopK,
NumCtx:        aiCfg.ContextLen,
RepeatPenalty: aiCfg.RepeatPen,
Seed:          aiCfg.Seed,
NumPredict:    aiCfg.MaxTokens,
},
}
jsonData, err := json.Marshal(reqPayload)
if err != nil {
pushLog("JSON Marshal Error in token stack: " + err.Error())
return
}
ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
defer cancel()
req, _ := http.NewRequestWithContext(ctx, "POST", aiCfg.OllamaAPI, bytes.NewBuffer(jsonData))
req.Header.Set("Content-Type", "application/json")
client := &http.Client{}
resp, err := client.Do(req)
if err != nil {
pushLog("OLLAMA EXECUTION ERROR: Core refused connection. Confirm port 11434 binding.")
return
}
defer resp.Body.Close()
bodyBytes, _ := io.ReadAll(resp.Body)
if resp.StatusCode != http.StatusOK {
pushLog(fmt.Sprintf("OLLAMA ERROR STATUS %d: %s", resp.StatusCode, string(bodyBytes)))
return
}
var ollamaResp OllamaChatResponse
if err := json.Unmarshal(bodyBytes, &ollamaResp); err != nil {
pushLog("Failed to decode token matrix JSON response: " + err.Error())
return
}
sessionHistories[currentSessionID] = append(sessionHistories[currentSessionID], OllamaMessage{
Role:    "assistant",
Content: ollamaResp.Message.Content,
})
pushLog(">>> REAL LLM FRAME DEPLOYMENT GENERATION COMPLETE.")
fyne.DoQueue(refreshChatUI)
}()
})
firewallBox := container.NewVBox(
firewallLabel,
container.NewGridWithColumns(2, allowBtn, blockBtn),
fwStatusLabel,
)
coreBox := container.NewVBox(
widget.NewLabel("## [ THE MYTHOS CORE ]"),
modelCheck,
executeModelBtn,
openParamBtn,
exportSecureLogBtn,
)
radarBox := container.NewVBox(
widget.NewLabel("## [ SPIDER & ONION RADAR ]"),
widget.NewLabel("AUTOMATED .ONION MONITORING TARGET:"),
darkWebTargetEntry,
widget.NewSeparator(),
widget.NewLabel("SURFACE WEB SEARCH:"),
surfaceWebQuery,
launchSurfaceBtn,
widget.NewLabel("INTERNAL PROXY CHANNEL:\nsocks5h://"+aiCfg.TorProxyAddr),
widget.NewLabel(">>> FIREWALL SHIELD REINFORCED. SYSTEMS ACTIVE."),
)
leftPanel := container.NewVBox(
firewallBox,
widget.NewSeparator(),
coreBox,
widget.NewSeparator(),
radarBox,
)
contextControlBar := container.NewGridWithColumns(3, sessionSelect, newSessionBtn, purgeSessionBtn)
sandboxBox := container.NewBorder(
container.NewVBox(widget.NewLabel("## [ CODE SANDBOX OS ]\n捕獲 AI DEPLOY CODE"), contextControlBar),
container.NewGridWithColumns(2, chatInput, sendPromptBtn),
nil, nil,
chatOutput,
)
mainLayout := container.NewGridWithColumns(3,
container.NewVScroll(leftPanel),
sandboxBox,
container.NewBorder(widget.NewLabel("== DEEPLOG TERMINAL =="), nil, nil, nil, logTerminal),
)
window.SetContent(mainLayout)
window.ShowAndRun()
}