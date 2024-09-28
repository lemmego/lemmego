import React from "react";
import { useForm, usePage } from "@inertiajs/react";

const Post: React.FC = () => {
  const { errors, input, message } = usePage().props;
  const { data, setData, post, progress } = useForm({
    title: "",
    body: "",
    logo: undefined,
    publish_at: "",
  });

  function handleInput(
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>
  ) {
    if ("files" in e.target && e.target.files && e.target.files.length > 0) {
      setData(e.target.name, e.target.files[0]);
      return;
    }

    setData(e.target.name, e.target.value);
  }

  function submit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    post("/post");
  }

	return (
		<div className="w-1/3 mx-auto">
			<h1 className="text-3xl text-center">Post</h1>
      {message && <p className="text-blue-500 text-center">{message}</p>}
			<form onSubmit={submit}>
				
					<div className="mt-2">
						<label htmlFor="title" className="label-primary">Title</label>
                        <input id="title" name="title" type="text" className="input" value={data.title} onChange={handleInput}/>
						{errors.title && <p className="text-xs text-red-500">{errors.title.join(', ')}</p>}
					</div>
				
					<div className="mt-2">
						<label htmlFor="body" className="label-primary">Body</label>
						<textarea id="body" name="body" className="input" value={data.body} onChange={handleInput}></textarea>
						{errors.body && <p className="text-xs text-red-500">{errors.body.join(', ')}</p>}
					</div>
				
					<div className="mt-2">
						<label htmlFor="logo" className="label-primary">Logo</label>
						<input id="logo" name="logo" type="file" onChange={handleInput}/>
						{errors.logo && <p className="text-xs text-red-500">{errors.logo.join(', ')}</p>}
					</div>
				
					<div className="mt-2">
						<label htmlFor="publish_at" className="label-primary">Publish At</label>
						<input id="publish_at" name="publish_at" type="date" className="input" value={data.publish_at} onChange={handleInput}/>
						{errors.publish_at && <p className="text-xs text-red-500">{errors.publish_at.join(', ')}</p>}
					</div>
				
				<div>
					<button type="submit" className="mt-4 btn-primary">Submit</button>
				</div>
			</form>
		</div>
	);
};

export default Post;

