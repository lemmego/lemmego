import { usePage, Link, Head, router } from "@inertiajs/react";
import { useState } from "react";

const TaskIndex = () => {
  const { props } = usePage();
  const tasks = (props.tasks || []) as any[];

  return (
    <div className="p-8 max-w-4xl mx-auto">
      <Head title="Tasks" />
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold">Tasks</h1>
        <Link href="/tasks/create"><button className="bg-blue-500 text-white px-4 py-2 rounded">New Task</button></Link>
      </div>
      {tasks.length === 0 ? (
        <p className="text-gray-500">No tasks yet.</p>
      ) : (
        <ul className="space-y-2">
          {tasks.map((t: any) => (
            <li key={t.id} className="flex items-center gap-3 p-3 border rounded">
              <Link href={`/tasks/${t.id}/edit`} className="flex-1 hover:underline">{t.title}</Link>
              <span className="text-sm text-gray-500">{t.status}</span>
              <button className="text-red-500 text-sm" onClick={() => { if (confirm("Delete?")) router.delete(`/tasks/${t.id}`); }}>Delete</button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default TaskIndex;
