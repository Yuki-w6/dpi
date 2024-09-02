import Home from './pages/Home';

export const BASE_URL =
  import.meta.env.MODE === 'development' ? 'http://localhost:1000/api' : '/api';

function App() {
  return <Home />;
}

export default App;
