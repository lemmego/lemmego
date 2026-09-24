import { usePage, Link, Head, router } from "@inertiajs/react";
import { useState } from "react";

const TaskForm = () => {
  const { props } = usePage();
  const task = props.task as any;
  const errors = (props.errors || {}) as Record<string, string>;
  const isEditing = !!task;
  const [title, setTitle] = useState(isEditing ? task.title : "");

  const submit = (e: any) => {
    e.preventDefault();
    if (isEditing) router.put(`/tasks/${task.id}`, { title });
    else router.post("/tasks", { title });
  };

  return (
    <div className="p-8 max-w-xl mx-auto">
      <Head title={isEditing ? "Edit Task" : "New Task"} />
      <h1 className="text-2xl font-bold mb-6">{isEditing ? "Edit Task" : "New Task"}</h1>
      <form onSubmit={submit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-1">Title</label>
          <input value={title} onChange={e => setTitle(e.target.value)} className="w-full border rounded px-3 py-2 text-sm" />
          {errors.title && <p className="text-red-500 text-xs mt-1">{errors.title}</p>}
        </div>
        <div className="flex gap-3">
          <button type="submit" className="bg-blue-500 text-white px-4 py-2 rounded text-sm">{isEditing ? "Update" : "Create"}</button>
          <Link href="/tasks"><button type="button" className="border px-4 py-2 rounded text-sm">Cancel</button></Link>
        </div>
      </form>
    </div>
  );
};

export default TaskForm;
