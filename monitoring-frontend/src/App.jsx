import { useEffect, useState } from 'react';
import { Line } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import api from './api';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

function App() {
  const [dataPoints, setDataPoints] = useState([]);

  useEffect(() => {
    // เรียก metrics ล่าสุด 1 ชั่วโมง
    api.get('/metrics?service=auth-service&name=require_latency_ms')
      .then(res => setDataPoints(res.data))
      .catch(console.error);
  }, []);

  const chartData = {
    labels: dataPoints.map(m => new Date(m.observedat).toLocaleTimeString()),
    datasets: [
      {
        label: 'Latency ms',
        data: dataPoints.map(m => m.value),
        borderColor: 'rgb(75, 192, 192)',
        tension: 0.2
      }
    ]
  };

  return (
    <div style={{ maxWidth: 800, margin: '40px auto' }}>
      <h1>Monitoring Dashboard</h1>
      <Line data={chartData} />
    </div>
  );
}

export default App;
