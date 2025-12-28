import { Link } from 'react-router-dom';

export default function Home() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center gap-4">
      <h1 className="text-3xl font-bold">Home</h1>

      <Link to="/register" className="text-blue-600 underline">
        Go to Register
      </Link>
    </div>
  );
}
