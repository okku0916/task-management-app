// frontend/src/app/page.tsx

import TaskForm from "./components/TaskForm";
import DeleteButton from "./components/DeleteButton";
import StatusButton from "./components/StatusButton";
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
			<div key={task.id} className={`p-5 border rounded-xl shadow-sm bg-white ${task.status === 'DONE' ? 'opacity-60' : ''}`}>
			<div className="flex justify-between items-start">
			<div>
			<h3 className={`font-bold text-lg ${task.status === 'DONE' ? 'line-through text-gray-500' : 'text-gray-800'}`}>
			{task.title}
			</h3>
			<p className="text-gray-600 mt-2 text-sm">{task.description}</p>
			</div>
			<div className="flex flex-col items-end gap-2">
			<span className={`px-2 py-1 text-[10px] font-black rounded uppercase ${task.status === 'DONE' ? 'bg-gray-200 text-gray-600' : 'bg-green-100 text-green-700'}`}>
			{task.status}
			</span>

			{/* ボタン類をまとめる */}
			<div className="flex gap-2 items-center">
			<StatusButton id={task.id} status={task.status} />
			<DeleteButton id={task.id} />
			</div>
			</div>
			</div>
			{/* ... 略 ... */}
			</div>
		))}
		</div>
		)}
		</section>
		</main>
	);
}
