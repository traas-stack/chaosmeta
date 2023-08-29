import * as echarts from 'echarts';
import { useEffect, useRef } from 'react';
import { ObservationChartContainer } from './style';

const ObservationCharts = () => {
  const chartRef = useRef<any>(null);
  const option = {
    title: {
      text: 'CPU',
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
      },
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: ['00:00', '01:15', '02:30', '03:45', '05:00', '06:15', '07:30', '08:45', '10:00', '11:15', '12:30', '13:45', '15:00', '16:15', '17:30', '18:45', '20:00', '21:15', '22:30', '23:45'],
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        formatter: '{value} W',
      },
      axisPointer: {
        snap: true,
      },
    },
    series: [
      {
        // name: 'Electricity',
        type: 'line',
        smooth: true,
        data: [300, 280, 250, 260, 270, 300, 550, 500, 400, 390, 380, 390, 400, 500, 600, 750, 800, 700, 600, 400],
        areaStyle: {},
        markLine: {
          symbol: ['none', 'none'],
          label: { show: false },
          data: [{ xAxis: '03:45' },{ xAxis: '06:15' }]
        },
      },
    ],
  };

  useEffect(() => {
    if (chartRef.current) {
      const currentChart = echarts.init(chartRef.current);
      if (option) {
        currentChart.setOption(option);
      }
    }
  }, []);

  return <ObservationChartContainer ref={chartRef} />;
};

export default ObservationCharts;
