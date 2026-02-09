// frontend/src/app/page.tsx

type Task = {
  id: number;
  title: string;
  description: string;
  status: string;
};

export default async function Home() {
  // 1. Dockerコンテナ間通信のため、localhostではなくサービス名「backend」を指定
  const res = await fetch("http://backend:8080/tasks", { cache: "no-store" });
  
  if (!res.ok) {
    throw new Error('Failed to fetch tasks');
  }

  const tasks: Task[] = await res.json();

  return (
    <main className="p-8">
      <h1 className="text-2xl font-bold mb-6">タスク一覧</h1>
      <div className="grid gap-4">
        {tasks.map((task: Task) => (
          <div key={task.id} className="p-4 border rounded shadow-sm bg-white">
            <h2 className="font-semibold text-lg">{task.title}</h2>
            <p className="text-gray-600">{task.description}</p>
            <span className="mt-2 inline-block px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded">
              {task.status}
            </span>
          </div>
        ))}
      </div>
    </main>
  );
}
