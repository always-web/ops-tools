var status_option = {
    tooltip : {
        trigger: 'item',
        formatter: "{a} <br/>{b} : {c} ({d}%)"
    },
    toolbox: {
        show : false
    },
    calculable : true,
    series : [
        {
            name:'状态码状态分布',
            type:'pie',
            radius : '55%',
            center: ['50%', '50%'],
            data:[]
        }
    ]
};

var hits_option = {
    tooltip : {
        trigger: 'axis'
    },
    toolbox: {
        show : false
    },
    calculable : true,
    xAxis : [
        {
            type : 'category',
            boundaryGap : false,
            data : []
        }
    ],
    yAxis : [
        {
            type : 'value'
        }
    ],
    series : [
        {
            name:'访问量',
            type:'line',
            stack: '访问量',
            data:[]
        }
    ]
};


var vistors_option = {
    tooltip : {
        trigger: 'axis',
        axisPointer : {
            type : 'shadow'
        }
    },
    toolbox: {
        show : false
    },
    calculable : true,
    xAxis : [
        {
            type : 'category',
            data : []
        }
    ],
    yAxis : [
        {
            type : 'value'
        }
    ],
    series : [
        {
            name:'用户数',
            type:'bar',
            data:[]
        }
    ]
};

var filesizeformat = function(bytes) {
    var units = ["B", "KB", "MB", "GB", "TB", "PB", "EB"];
    var idx = 0;
    while(bytes >= 1000) {
        bytes /= 1000;
        idx ++;
    }

    return bytes.toFixed(2) + units[idx];
}

var bytes_option = {
    tooltip : {
        trigger: 'axis',
        formatter: function(a) {
            return a[0].name + "<br/>流量:" + filesizeformat(a[0].data);
        }
    },
    toolbox: {
        show : false
    },
    calculable : true,
    xAxis : [
        {
            type : 'category',
            boundaryGap : false,
            data : []
        }
    ],
    yAxis : [
        {
            type : 'value'
        }
    ],
    series : [
        {
            name:'流量',
            type:'line',
            stack: '流量',
            itemStyle: {normal: {areaStyle: {type: 'default'}}},
            data:[]
        }
    ]
};



var region_option = {
    backgroundColor: '#1b1b1b',
    color: ['gold','aqua','lime'],
    tooltip : {
        trigger: 'item',
        formatter: '{b}'
    },
    toolbox: {
        show : false
    },
    dataRange: {
        min : 0,
        max : 100,
        calculable : true,
        color: ['#ff3333', 'orange', 'yellow','lime','aqua'],
        textStyle:{
            color:'#fff'
        }
    },
    series : [
        {
            name: '全国',
            type: 'map',
            roam: true,
            hoverable: false,
            mapType: 'china',
            itemStyle:{
                normal:{
                    borderColor:'rgba(100,149,237,1)',
                    borderWidth:0.5,
                    areaStyle:{
                        color: '#1b1b1b'
                    }
                }
            },
            data:[],
            geoCoord: {}
        },
        {
            name: '北京',
            type: 'map',
            mapType: 'china',
            data:[],
            markLine : {
                smooth:true,
                effect : {
                    show: true,
                    scaleSize: 1,
                    period: 30,
                    color: '#fff',
                    shadowBlur: 10
                },
                itemStyle : {
                    normal: {
                        borderWidth:1,
                        lineStyle: {
                            type: 'solid',
                            shadowBlur: 10
                        }
                    }
                },
                data : []
            },
            markPoint : {
                symbol:'emptyCircle',
                symbolSize : function (v){
                    return 10 + v/10000
                },
                effect : {
                    show: true,
                    shadowBlur : 0
                },
                itemStyle:{
                    normal:{
                        label:{show:false}
                    },
                    emphasis: {
                        label:{position:'top'}
                    }
                },
                data : []
            }
        }
    ]
};
