import ProjectList from "../../components/project/ProjectList";

export default function ProjectsPage() {
  return (
    <div className="max-w-6xl mx-auto p-4">
      <div className="flex items-center justify-between mb-4">
        <h1 className="text-2xl font-bold">My Projects</h1>
      </div>
      <ProjectList />
    </div>
  );
}
