package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/oschwald/geoip2-golang"
	"github.com/otiai10/copy"
	"html/template"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	// daysMap 定义日期
	daysMap = make(map[string]bool)
	// hitTotal 总的点击量
	hitTotal int
	// bytesTotal 总的网络流量
	bytesTotal uint64
	// ipTotalMap 每个 ip 的访问次数
	ipTotalMap = make(map[string]int)
	// statusTotalMap 每个 http 状态码出现的次数
	statusTotalMap = make(map[string]int)
	// hitDaysMap 每天的访问量
	hitDaysMap = make(map[string]int)
	// bytesDaysMap 每天的流量大小
	bytesDaysMap = make(map[string]uint64)
	// ipDaysMap 每天每个 ip 的访问量
	ipDaysMap = make(map[string]map[string]int)
	// statusDaysMap 每个 http 状态码每天的访问量
	statusDaysMap = make(map[string]map[string]int)
)

func main() {
	path := flag.String("path", "", "nginx 日志路径")
	dir := flag.String("dir", "", "报告输出地址")
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", "loganalysis")
		flag.PrintDefaults()
	}
	flag.Parse()
	fmt.Println("path:", *path)
	fmt.Println("dir:", *dir)

	if *dir == "" {
		*dir = fmt.Sprintf("reports/report_%d", time.Now().Unix())
	}
	if _, err := os.Stat(*dir); err == nil {
		fmt.Printf("结果目录 %s 已存在\n", *dir)
		os.Exit(-1)
	}

	file, err := os.Open(*path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("路径 %s 不存在\n", *path)
		} else {
			fmt.Printf("无法打开 nginx 日志,err: %v\n", err)
		}
		os.Exit(-1)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	if fileInfo, _ := file.Stat(); fileInfo.IsDir() {
		fmt.Printf("路径 %s 不能是一个目录\n", *path)
		_ = file.Close()
		os.Exit(-1)
	}
	if err := processNginxLog(file, dir); err != nil {
		fmt.Println("处理 nginx 日志 error: ", err)
		_ = file.Close()
		os.Exit(-1)
	}

	if err := handleTemplate(dir); err != nil {
		fmt.Println("渲染 template error: ", err)
		_ = file.Close()
		os.Exit(-1)
	}

}

func handleTemplate(dir *string) error {
	funcMaps := template.FuncMap{
		"filesizeformat": humanize.Bytes,
		"sortmap":        SortMap,
		"json":           Json,
	}
	tmpl, err := template.New("index.html").Funcs(funcMaps).ParseFiles(filepath.Join("tpl", "index.html"))
	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(*dir, "index.html"))
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	writer := bufio.NewWriter(file)
	defer func(writer *bufio.Writer) {
		_ = writer.Flush()
	}(writer)

	dayKeys := make([]string, len(daysMap))
	index := 0
	for key := range daysMap {
		dayKeys[index] = key
		index++
	}

	geoip, err := geoip2.Open("db/GeoLite2-City.mmdb")
	if err != nil {
		return err
	}
	defer func(geoip *geoip2.Reader) {
		_ = geoip.Close()
	}(geoip)

	// 地址和访问次数的映射
	regionTotalMap := make(map[string]int)
	// 地址的经纬度
	regionLocationMap := make(map[string][2]float64)
	for ipStr, count := range ipTotalMap {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			fmt.Println("ipStr can't convert to ip", ipStr)
			continue
		}
		record, err := geoip.City(ip)
		if err != nil {
			fmt.Printf("ipStr %s can't convert to city,err:%v\n", ipStr, err)
			continue
		}

		// 只显示国内 IP 地址
		if record.Country.Names["zh-CN"] != "中国" {
			continue
		}

		name := fmt.Sprintf("%s%s", record.Country.Names["zh-CN"], record.City.Names["zh-CN"])
		regionTotalMap[name] += count
		if _, exists := regionLocationMap[name]; !exists {
			regionLocationMap[name] = [2]float64{record.Location.Longitude, record.Location.Latitude}
		}
	}

	err = tmpl.Execute(writer, struct {
		Days           []string
		HitTotal       int
		BytesTotal     uint64
		VistorsTotal   IntMapOrdered
		StatusTotal    map[string]int
		HitDays        map[string]int
		BytesDays      map[string]uint64
		VistorsDays    map[string]map[string]int
		StatusDays     map[string]map[string]int
		RegionTotal    map[string]int
		RegionLocation map[string][2]float64
	}{
		Days:           dayKeys,
		HitTotal:       hitTotal,
		BytesTotal:     bytesTotal,
		VistorsTotal:   NewIntMapOrdered(ipTotalMap),
		StatusTotal:    statusTotalMap,
		HitDays:        hitDaysMap,
		BytesDays:      bytesDaysMap,
		VistorsDays:    ipDaysMap,
		StatusDays:     statusDaysMap,
		RegionTotal:    regionTotalMap,
		RegionLocation: regionLocationMap,
	})
	if err != nil {
		return err
	}
	return nil
}

func processNginxLog(file *os.File, dir *string) error {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		nodes := strings.Split(line, " ")
		if len(nodes) < 12 {
			continue
		}
		// nginx 日志格式是: ip datetime method url status bytes
		logTime, _ := time.Parse("02/Jan/2006:15:04:05", nodes[3][1:])
		logDay := logTime.Format("2006-01-02")

		daysMap[logDay] = true

		hitTotal++
		hitDaysMap[logDay]++

		if b, err := strconv.ParseUint(nodes[9], 10, 64); err == nil {
			bytesTotal += b
			bytesDaysMap[logDay] += b
		}

		if _, exists := ipDaysMap[logDay]; !exists {
			ipDaysMap[logDay] = make(map[string]int)
		}

		if _, exists := statusDaysMap[logDay]; !exists {
			statusDaysMap[logDay] = make(map[string]int)
		}

		ipTotalMap[nodes[0]]++
		ipDaysMap[logDay][nodes[0]]++

		statusTotalMap[nodes[8]]++
		statusDaysMap[logDay][nodes[8]]++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := copy.Copy("tpl", *dir); err != nil {
		return err
	}
	return nil
}
