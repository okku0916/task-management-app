"use client";

import { useRouter } from "next/navigation";

export default function StatusButton({ id, status }: { id: number; status: string }) {
	const router = useRouter();

	const handleUpdate = async () => {
		const res = await fetch(`http://localhost:8080/tasks?id=${id}`, {
			method: "PUT",
		});

		if (res.ok) {
			router.refresh();
		}
	};

	if (status === "DONE") return null; // すでに完了ならボタンを出さない

	return (
		<button 
		onClick={handleUpdate}
		className="text-xs bg-green-500 hover:bg-green-600 text-white font-bold py-1 px-2 rounded transition-colors"
		>
		完了にする
		</button>
	);
}
