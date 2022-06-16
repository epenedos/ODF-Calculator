package main
import (
//    "os"
    "fmt"
    "net/http"
    "net/url"
    "log"
    "html/template"
    "strconv"
//    "math/rand"
    "github.com/go-echarts/go-echarts/v2/charts"
    "github.com/go-echarts/go-echarts/v2/opts"
//    "github.com/go-echarts/go-echarts/v2/types"
)


      var number_nodes int
      var number_cores int
      var disk_size float32
      var number_disks_node  int
      var output_type int = 1
      var html string



func calculate(number_nodes int,number_cores int, disk_size float32,number_disks_node int,w http.ResponseWriter) {
      var master_cpu int = 10
      var master_mem  int = 12
      var disk_cpu  int = 2
      var disk_mem int = 5

      var node_cpu int
      var node_mem  int
      var node_space float32


      var cluster_cpu int
      var cluster_mem  int
      var cluster_space float32

      w.Write([]byte("<div class=\"results\">"))
      w.Write([]byte("<div class=\"cluster-list\">"))
      pos:= 1
      for i := 0; i < number_nodes; i++ {
        node_cpu=0
        node_mem=0
        node_space=0
        if pos == 1 {
          w.Write([]byte("<div class=\"leftdiv\">"))
        }
        if pos == 2 {
          w.Write([]byte("<div class=\"middlediv\">"))
        }
        if pos == 3 {
          w.Write([]byte("<div class=\"rightdiv\">"))
        }
        if pos < 3 {
          pos = pos + 1
        } else {pos = 1}

          if i<3 {
              node_cpu=master_cpu
              node_mem=master_mem
              html = fmt.Sprint("<p><img src='/assets/img/fallbackhardware.svg' width='20' height='20'><b>  Master Node: ",i+1, "</b> CPU ODF: ",master_cpu," cpu's", " Mem ODF: ",master_mem, " Gb</p>")
          } else {
              html = fmt.Sprint("<p><img src='/assets/img/fallbackhardware.svg' width='20' height='20'><b>  Worker Node: ",i+1," </b> </p>")
          }
          w.Write([]byte(html))



          for j := 0;  j < number_disks_node; j++ {
              node_cpu=node_cpu+disk_cpu
              node_mem=node_mem+disk_mem
              node_space=node_space+disk_size
              html = fmt.Sprint("<p><img src='/assets/img/Icon-Red_Hat-Storage-A-Red-RGB.svg' width='20' height='20'>  Disk: ",j+1," CPU ODS: ",disk_cpu ," cpu's"," Mem ODS: ",disk_mem ," Gb </p>")
              w.Write([]byte(html))
          }

          html = fmt.Sprint("<p><b>Total CPU of node : ",i+1,"</b> total with HT ", number_cores * 2 , " allocated to ODF", node_cpu," cpu's</p>")
          w.Write([]byte(html))
          html = fmt.Sprint("<p><b>Total Mem of node : ",i+1,"</b> # ", node_mem, " Gb</p>")
          w.Write([]byte(html))
          html = fmt.Sprint("<p><img src='/assets/img/Icon-Red_Hat-Storage_stack-A-Red-RGB.svg' width='20' height='20'><b>Total Storage Space of node : ",i+1,"</b> # ", node_space, " Tb (RAW)</p>")
          w.Write([]byte(html))
          w.Write([]byte("</div>"))
          cluster_cpu=cluster_cpu+node_cpu
          cluster_mem=cluster_mem+node_mem
          cluster_space=cluster_space+node_space


    }
    w.Write([]byte("</div>"))
    w.Write([]byte("<div class=\"cluster-totals\">"))
    w.Write([]byte("<br><br><p><div class='rightdiv' ><b> Cluster totals </b></p>"))
    html = fmt.Sprint("<img src='/assets/img/Icon-Red_Hat-Datacenter-A-Red-RGB.svg' width='20' height='20'><p>Total CPU of cluster :  # ", cluster_cpu, " cpu <br>")
    w.Write([]byte(html))
    html = fmt.Sprint("Total Mem of cluster :  # ", cluster_mem, " Gb<br>")
    w.Write([]byte(html))
    html = fmt.Sprint("<img src='/assets/img/Icon-Red_Hat-Storage_stack-A-Red-RGB.svg' width='20' height='20'> <b>Total Storage Space of cluster :  # </b>", cluster_space, " Tb (RAW)<br>")
    w.Write([]byte(html))
    html = fmt.Sprint("<img src='/assets/img/Icon-Red_Hat-Storage_stack-A-Red-RGB.svg' width='20' height='20'> <b>Total Storage Space of cluster :  # </b>", cluster_space/3, " Tb (Usable Capacity)<br></div></p>")
    w.Write([]byte(html))
    w.Write([]byte("</div>"))
    w.Write([]byte("</div>"))
}

func main(){

  fs := http.FileServer(http.Dir("assets"))
  mux := http.NewServeMux()
  mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
  mux.HandleFunc("/graph", httpserver)
  mux.HandleFunc("/", httpserver_home)
  fmt.Println("Server started at port 8081")
  log.Fatal(http.ListenAndServe(":8081", mux))

}

// generate random data for bar chart
func generateBarItems(raw bool) []opts.BarData {
    items := make([]opts.BarData, 0)
    for i := 0; i < 6; i++ {
      if raw {
        items = append(items, opts.BarData{Value: float32(i+1)*(disk_size*float32(number_nodes))})
      } else {
        items = append(items, opts.BarData{Value: float32(i+1)*(disk_size*float32(number_nodes)/3)})
      }
    }
    return items
}



func httpserver_home(w http.ResponseWriter, r *http.Request) {
    var tpl = template.Must(template.ParseFiles("index.html"))
    tpl.Execute(w, nil)

}


func httpserver(w http.ResponseWriter, r *http.Request) {

  u, err := url.Parse(r.URL.String())
  	if err != nil {
  		http.Error(w, err.Error(), http.StatusInternalServerError)
  		return
  	}

  	params := u.Query()
  	nn, err := strconv.ParseInt(params.Get("nhosts"), 10, 0)
    number_nodes=int(nn)
    nc,err := strconv.ParseInt(params.Get("ncores"),10,0)
    number_cores = int(nc)
    ndn, err := strconv.ParseInt(params.Get("ndisks"),10,0)
    number_disks_node=int(ndn)
    ds,err := strconv.ParseFloat(params.Get("dsize"),0)
    disk_size=float32(ds)

    html = fmt.Sprint("<html lang='en'>")
    w.Write([]byte(html))
    html = fmt.Sprint("<head><title>ODF Storage Calculator </title> <link rel='stylesheet' href='/assets/style.css' /></head>")
    w.Write([]byte(html))

    html = fmt.Sprint("<body>")
    w.Write([]byte(html))

    html = fmt.Sprint("<div class='fixed-header'>")
    w.Write([]byte(html))

    html = fmt.Sprint("<img id='logo' src='/assets/img/Product_Icon-Red_Hat-OpenShift-RGB.svg' width='75' height='75'> <span id='title'>OpenShift Data Foundation Calculator</span>")
    w.Write([]byte(html))
    html = fmt.Sprint("</div>")
    w.Write([]byte(html))

    calculate(number_nodes,number_cores,disk_size,number_disks_node, w)


  bar := charts.NewBar()

  // Set global options
  bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
      Title:    "Storage Calculator",
      Subtitle: fmt.Sprintf("%v", disk_size) + " Disks per node with RF=2 (3 copies of the data, Tolerates 2 node failure)",
  }))

  // Put data into instance
  bar.SetXAxis([]string{"1", "2", "3", "4", "5", "6"}).
      AddSeries("Raw", generateBarItems(true)).
      AddSeries("Usable Capacity", generateBarItems(false))
  bar.Render(w)



  html = fmt.Sprint("<div class='fixed-footer'>KAM Software Solutions</div>")
  w.Write([]byte(html))
  html = fmt.Sprint("</body>")
  w.Write([]byte(html))
  html = fmt.Sprint("</html>")
  w.Write([]byte(html))




}
