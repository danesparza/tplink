package tplink

import (
	"encoding/json"
	"fmt"
	"time"
)

type HS100 struct {
	ip string
}

// Get System Info (Software & Hardware Versions, MAC, deviceID, hwID etc.)
func (p *HS100) Info() (*Info, error) {
	data, err := exec(p.ip, GET_INFO)
	if err != nil {
		return nil, err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return nil, err
	}

	return r.System.Info, nil
}

// Reboot
func (p *HS100) Reboot() (string, error) {
	return exec(p.ip, REBOOT)
}

// Reset
func (p *HS100) Reset() (string, error) {
	return exec(p.ip, RESET)
}

// Set alias/name
func (p *HS100) SetAlias(alias string) error {
	data, err := exec(p.ip, fmt.Sprintf(SET_ALIAS, alias))
	if err != nil {
		return err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return err
	}

	if r.System.SetAlias.ErrorCode != 0 {
		return fmt.Errorf("failed to set alias. Error code=%d", r.System.SetAlias.ErrorCode)
	}
	return nil
}

// Turn On
func (p *HS100) TurnOn() error {
	data, err := exec(p.ip, TURN_ON)
	if err != nil {
		return err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return err
	}

	if r.System.SetState.ErrorCode != 0 {
		return fmt.Errorf("failed to turn the device off. Error code=%d", r.System.SetState.ErrorCode)
	}
	return nil
}

// Turn Off
func (p *HS100) TurnOff() error {
	data, err := exec(p.ip, TURN_OFF)
	if err != nil {
		return err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return err
	}

	if r.System.SetState.ErrorCode != 0 {
		return fmt.Errorf("failed to turn the device off. Error code=%d", r.System.SetState.ErrorCode)
	}
	return nil
}

// Turn Led Light On
func (p *HS100) TurnLedOn() error {
	data, err := exec(p.ip, TURN_LED_ON)
	if err != nil {
		return err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return err
	}

	if r.System.SetState.ErrorCode != 0 {
		return fmt.Errorf("failed to turn the device off. Error code=%d", r.System.SetState.ErrorCode)
	}
	return nil
}

// Turn Led Light Off
func (p *HS100) TurnLedOff() error {
	data, err := exec(p.ip, TURN_LED_OFF)
	if err != nil {
		return err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return err
	}

	if r.System.SetState.ErrorCode != 0 {
		return fmt.Errorf("failed to turn the device off. Error code=%d", r.System.SetState.ErrorCode)
	}
	return nil
}

// TODO: return a timezone instead of index
func (p *HS100) TimeZone() (int, error) {
	data, err := exec(p.ip, GET_TIMEZONE)
	if err != nil {
		return 0, err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return 0, err
	}

	tz := r.Time.GetTimeZone

	if tz.ErrorCode != 0 {
		return 0, fmt.Errorf("failed to get time zone. Error code=%d", r.System.SetState.ErrorCode)
	}

	return tz.Index, nil
}

func (p *HS100) Time() (time.Time, error) {
	data, err := exec(p.ip, GET_TIME)
	if err != nil {
		return time.Time{}, err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return time.Time{}, err
	}

	t := r.Time.GetTime

	if t.ErrorCode != 0 {
		return time.Time{}, fmt.Errorf("failed to get time. Error code=%d", r.System.SetState.ErrorCode)
	}

	// TODO: get timezone
	//timezone, err := p.TimeZone()
	// if err != nill {
	//	return "", fmt.Errorf("failed to get device timezone: %s", err)
	//}

	loc, err := time.LoadLocation("EST")
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get device timezone: %s", err)
	}
	d := time.Date(t.Year, time.Month(t.Month), t.Day, t.Hour, t.Minutes, t.Seconds, 0, loc)

	return d, nil
}

func (p *HS100) SetTimeZone(t time.Time) error {
	// TODO: timezone
	cmd := fmt.Sprintf(SET_TIMEZONE, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 18)
	data, err := exec(p.ip, cmd)
	if err != nil {
		return err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return err
	}

	if r.Time.SetTimeZone.ErrorCode != 0 {
		return fmt.Errorf("failed to set timezone. Error code=%d, msg: %s", r.Time.SetTimeZone.ErrorCode, r.Time.SetTimeZone.ErrorMessage)
	}

	return nil
}

func (p *HS100) ScanWifi() ([]AP, error) {
	data, err := exec(p.ip, SCAN_WIFI)
	if err != nil {
		return nil, err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return nil, err
	}

	if r.NetIf.GetScanInfo.ErrorCode != 0 {
		return nil, fmt.Errorf("failed to scan for wifi networks. Error code=%d, msg: %s", r.NetIf.GetScanInfo.ErrorCode, r.NetIf.GetScanInfo.ErrorMessage)
	}

	return r.NetIf.GetScanInfo.List, nil
}

func (p *HS100) SetWifi(ssid string, password string, keyType int) error {
	cmd := fmt.Sprintf(SET_WIFI, ssid, password, keyType)
	data, err := exec(p.ip, cmd)
	if err != nil {
		return err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return err
	}

	if r.NetIf.SetWifi.ErrorCode != 0 {
		return fmt.Errorf("failed to set wifi. Error code=%d, msg: %s", r.NetIf.SetWifi.ErrorCode, r.NetIf.SetWifi.ErrorMessage)
	}

	return nil
}

// Gets Cloud Info (Server, Username, Connection Status)
func (p *HS100) CloudInfo() (*Cloud, error) {
	data, err := exec(p.ip, GET_CLOUD_INFO)
	if err != nil {
		return nil, err
	}

	fmt.Println(data)

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return nil, err
	}

	if r.CNCloud.Info.ErrorCode != 0 {
		return nil, fmt.Errorf("failed to get cloud info. Error code=%d, msg: %s", r.CNCloud.Info.ErrorCode, r.CNCloud.Info.ErrorMessage)
	}

	c := &Cloud{
		Username: r.CNCloud.Info.Username,
		Server:   r.CNCloud.Info.Server,
		Binded:   r.CNCloud.Info.Binded,
	}
	return c, nil
}

// Set Server URL
func (p *HS100) SetCloudUrl(url string) error {
	cmd := fmt.Sprintf(SET_CLOUD_URL, url)
	data, err := exec(p.ip, cmd)
	if err != nil {
		return err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return err
	}

	if r.CNCloud.SetServerUrl.ErrorCode != 0 {
		return fmt.Errorf("failed to get cloud info. Error code=%d, msg: %s", r.CNCloud.SetServerUrl.ErrorCode, r.CNCloud.SetServerUrl.ErrorMessage)
	}

	return nil
}

// Connects with server using username & Password
func (p *HS100) CloudBind(username string, password string) error {
	cmd := fmt.Sprintf(CLOUD_BIND, username, password)
	data, err := exec(p.ip, cmd)
	if err != nil {
		return err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return err
	}

	if r.CNCloud.Bind.ErrorCode != 0 {
		return fmt.Errorf("failed to bind to cloud. Error code=%d, msg: %s", r.CNCloud.Bind.ErrorCode, r.CNCloud.Bind.ErrorMessage)
	}

	return nil
}

// Unregister Device from Cloud Account
func (p *HS100) CloudUnbind() error {
	data, err := exec(p.ip, CLOUD_UNBIND)
	if err != nil {
		return err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return err
	}

	if r.CNCloud.Unbind.ErrorCode != 0 {
		return fmt.Errorf("failed to unbind devide from cloud. Error code=%d, msg: %s", r.CNCloud.Unbind.ErrorCode, r.CNCloud.Unbind.ErrorMessage)
	}

	return nil
}

// Gets Next Scheduled Action
func (p *HS100) GetNextScheduledAction() (*NextAction, error) {
	data, err := exec(p.ip, GET_NEXT_SCHEDULE_ACTION)
	if err != nil {
		return nil, err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return nil, err
	}

	if r.Schedule.GetNextAction.ErrorCode != 0 {
		return nil, fmt.Errorf("failed to get next scheduled action. Error code=%d, msg: %s", r.Schedule.GetNextAction.ErrorCode, r.Schedule.GetNextAction.ErrorMessage)
	}

	resp := &NextAction{
		RuleID:              r.Schedule.GetNextAction.RuleID,
		ScheduledTimeSecond: r.Schedule.GetNextAction.ScheduledTimeSecond,
		Action:              r.Schedule.GetNextAction.Action,
		Type:                r.Schedule.GetNextAction.Type,
	}
	return resp, nil
}

// Gets Schedule Rules List
func (p *HS100) GetScheduleList() ([]Rule, error) {
	data, err := exec(p.ip, GET_SCHEDULE_RULES_LIST)
	if err != nil {
		return nil, err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return nil, err
	}

	if r.Schedule.Rule.ErrorCode != 0 {
		return nil, fmt.Errorf("failed to get scheduled rules from device. Error code=%d, msg: %s", r.Schedule.Rule.ErrorCode, r.Schedule.Rule.ErrorMessage)
	}

	return r.Schedule.Rule.List, nil
}

// Add New Schedule Rule
func (p *HS100) AddScheduleRule(name string, days Days, action Action, minutes int, enable int, year int, month int, day int) (string, error) {
	return p.addScheduleRule(NONE, name, days, action, minutes, enable, year, month, day)
}

func (p *HS100) AddSunSetScheduleRule(name string, days Days, action Action, enable int, year int, month int, day int) (string, error) {
	return p.addScheduleRule(SUNSET, name, days, action, 0, enable, year, month, day)
}

func (p *HS100) AddSunRiseScheduleRule(name string, days Days, action Action, enable int, year int, month int, day int) (string, error) {
	return p.addScheduleRule(SUNRISE, name, days, action, 0, enable, year, month, day)
}

func (p *HS100) addScheduleRule(timeOpt TimeOption, name string, days Days, action Action, minutes int, enable int, year int, month int, day int) (string, error) {
	weekdays := days.String()
	repeat := OFF

	if weekdays != "[0,0,0,0,0,0,0]" {
		repeat = ON
	}
	cmd := fmt.Sprintf(ADD_SCHEDULE_RULE, timeOpt, weekdays, minutes, enable, repeat, name, month, action, year, day)
	data, err := exec(p.ip, cmd)
	if err != nil {
		return "", err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return "", err
	}

	if r.Schedule.AddRule.ErrorCode != 0 {
		return "", fmt.Errorf("failed to add scheduled rules. Error code=%d, msg: %s", r.Schedule.AddRule.ErrorCode, r.Schedule.AddRule.ErrorMessage)
	}

	return r.Schedule.AddRule.ID, nil
}

// Edit Schedule Rule with given ID
func (p *HS100) EditScheduleRule(id string, timeOpt TimeOption, name string, days Days, action Action, minutes int, enable int, year int, month int, day int) error {
	weekdays := days.String()
	repeat := OFF

	if weekdays != "[0,0,0,0,0,0,0]" {
		repeat = ON
	}
	cmd := fmt.Sprintf(EDIT_SCHEDULE_RULE, timeOpt, weekdays, minutes, enable, repeat, id, name, month, action, year, day)
	data, err := exec(p.ip, cmd)
	if err != nil {
		return err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return err
	}

	if r.Schedule.EditRule.ErrorCode != 0 {
		return fmt.Errorf("failed to edit scheduled rules. Error code=%d, msg: %s", r.Schedule.EditRule.ErrorCode, r.Schedule.EditRule.ErrorMessage)
	}

	return nil
}

// Delete Schedule Rule with given ID
func (p *HS100) DeleteScheduleRule(id string) error {
	cmd := fmt.Sprintf(DELETE_SCHEDULE_RULE, id)
	data, err := exec(p.ip, cmd)
	if err != nil {
		return err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return err
	}

	if r.Schedule.DeleteRule.ErrorCode != 0 {
		return fmt.Errorf("failed to edit scheduled rules. Error code=%d, msg: %s", r.Schedule.DeleteRule.ErrorCode, r.Schedule.DeleteRule.ErrorMessage)
	}

	return nil
}

// Delete All Schedule Rules and Erase Statistics
func (p *HS100) DeleteAllScheduleRule() error {
	data, err := exec(p.ip, DELETE_ALL_SCHEDULE_RULE)
	if err != nil {
		return err
	}

	r := Response{}
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return err
	}

	if r.Schedule.DeleteAllRules.ErrorCode != 0 {
		return fmt.Errorf("failed to edit scheduled rules. Error code=%d, msg: %s", r.Schedule.DeleteAllRules.ErrorCode, r.Schedule.DeleteAllRules.ErrorMessage)
	}

	return nil
}

func NewHS100(ip string) *HS100 {
	return &HS100{ip: ip}
}
