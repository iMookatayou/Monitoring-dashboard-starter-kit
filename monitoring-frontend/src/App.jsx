import { useState } from 'react';
import { Line } from 'react-chartjs-2';
import {
  Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend
} from 'chart.js';
import api from './api';

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend);

export default function App() {
  const [service, setService] = useState('auth-service');
  const [name, setName] = useState('require_latency_ms'); // ให้ตรงกับที่ POST
  const [hours, setHours] = useState(6);
  const [rows, setRows] = useState([]);
  const [loading, setLoading] = useState(false);
  const [err, setErr] = useState('');

  async function load() {
    setLoading(true);
    setErr('');
    try {
      const to = new Date();
      const from = new Date(Date.now() - Number(hours) * 60 * 60 * 1000);

      const res = await api.get('/metrics', {
        params: {
          service,
          name,
          from: from.toISOString(),
          to: to.toISOString()
        }
      });

      const data = Array.isArray(res.data) ? res.data : [];

      // รองรับคีย์ทั้ง snake_case และ PascalCase จาก backend
      const normalized = data.map(r => ({
        observed_at: r.observed_at ?? r.ObservedAt,
        value: r.value ?? r.Value ?? 0
      }));

      setRows(normalized.reverse()); // เรียงเวลาเก่า->ใหม่สำหรับกราฟ
    } catch (e) {
      setErr(e?.message || 'load error');
      setRows([]);
    } finally {
      setLoading(false);
    }
  }

  const chartData = {
    labels: rows.map(r => new Date(r.observed_at).toLocaleTimeString()),
    datasets: [
      {
        label: name || 'value',
        data: rows.map(r => r.value),
        borderColor: 'rgb(75, 192, 192)',
        backgroundColor: 'rgba(75, 192, 192, 0.2)',
        tension: 0.2,
        pointRadius: 2
      }
    ]
  };

  const chartOptions = {
    responsive: true,
    plugins: { legend: { position: 'top' }, title: { display: false } },
    animation: false
  };

  return (
    <div style={{ maxWidth: 1000, margin: '40px auto', color: '#eee' }}>
      <h1>Monitoring Dashboard</h1>

      <div style={{ display: 'flex', gap: 12, flexWrap: 'wrap', marginBottom: 16 }}>
        <div>
          <div style={{ fontSize: 12, opacity: .8 }}>Service</div>
          <input value={service} onChange={e => setService(e.target.value)} />
        </div>
        <div>
          <div style={{ fontSize: 12, opacity: .8 }}>Metric name</div>
          <input value={name} onChange={e => setName(e.target.value)} />
        </div>
        <div>
          <div style={{ fontSize: 12, opacity: .8 }}>Range (hours)</div>
          <select value={hours} onChange={e => setHours(e.target.value)}>
            <option value={1}>1</option>
            <option value={6}>6</option>
            <option value={12}>12</option>
            <option value={24}>24</option>
          </select>
        </div>
        <button onClick={load} disabled={loading}
          style={{ padding: '8px 14px', borderRadius: 8, border: '1px solid #444', background: '#1f6feb', color: '#fff' }}>
          {loading ? 'Loading...' : 'Load'}
        </button>
      </div>

      {err && <div style={{ color: '#ff6b6b', marginBottom: 12 }}>Error: {err}</div>}

      <div style={{ background: '#111', padding: 12, borderRadius: 12 }}>
        <Line data={chartData} options={chartOptions} />
      </div>
    </div>
  );
}
