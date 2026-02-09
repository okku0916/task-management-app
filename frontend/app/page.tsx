// frontend/src/app/page.tsx

import TaskForm from "./components/TaskForm";
type Task = {
  id: number;
  title: string;
  description: string;
  status: string;
  created_at: string;
};

export default async function Home() {
  // 1. Dockerコンテナ間通信のため、localhostではなくサービス名「backend」を指定
  const res = await fetch("http://backend:8080/tasks", { cache: "no-store" });
  
  if (!res.ok) {
    return <div className="p-8 text-red-500">APIサーバーとの接続に失敗しました。</div>;
  }
  const tasks: Task[] = await res.json();

  return (
    <main className="p-8 max-w-2xl mx-auto min-h-screen bg-gray-50">
      <header className="mb-10">
        <h1 className="text-4xl font-black text-gray-900 tracking-tight">Task Master</h1>
        <p className="text-gray-500 font-medium">Go & Next.js Fullstack Project</p>
      </header>

      {/* 登録フォームの表示 */}
      <TaskForm />

      <section>
        <h2 className="text-2xl font-bold mb-5 text-gray-800 border-b pb-2">現在のタスク</h2>
        {tasks.length === 0 ? (
          <p className="text-gray-400 italic">タスクはまだありません。上のフォームから追加しましょう！</p>
        ) : (
          <div className="grid gap-4">
            {tasks.map((task: Task) => (
              <div key={task.id} className="p-5 border border-gray-200 rounded-xl shadow-sm bg-white">
                <div className="flex justify-between items-start">
                  <h3 className="font-bold text-lg text-gray-800">{task.title}</h3>
                  <span className="px-2 py-1 bg-green-100 text-green-700 text-xs font-bold rounded uppercase">
                    {task.status}
                  </span>
                </div>
                <p className="text-gray-600 mt-2 text-sm">{task.description}</p>
                <time className="block mt-4 text-[10px] text-gray-400 uppercase tracking-widest">
                  {new Date(task.created_at).toLocaleString('ja-JP')}
                </time>
              </div>
            ))}
          </div>
        )}
      </section>
    </main>
  );
}
