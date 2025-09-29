import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE || 'http://localhost:8080',
  headers: {
    'X-API-Key': import.meta.env.VITE_API_KEY || 'supersecretapikey'
  }
});

export default api;
