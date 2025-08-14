import { Calendar, Eye, Package, Pencil, Trash2, EyeIcon } from "lucide-react";
import type { ProjectModel } from "@/api/manager/types";
import { Link } from "react-router-dom";

function Card({ children, className = "" }: { children: React.ReactNode; className?: string }) {
  return <div className={`border rounded-lg p-4 shadow-sm bg-white ${className}`}>{children}</div>;
}

export default function ProjectListCard({ project, onEdit, onDelete }: {
  project: ProjectModel;
  onEdit?: (project: ProjectModel) => void;
  onDelete?: (project: ProjectModel) => void;
}) {
  return (
    <Card>
      <div className="flex items-start justify-between">
        <div>
          <h3 className="text-lg font-semibold">{project.name}</h3>
          <p className="text-sm text-gray-500">/{project.slug}</p>
          {project.description && (
            <p className="text-sm text-gray-700 mt-1 line-clamp-2">{project.description}</p>
          )}
        </div>
        <div className="flex items-center gap-2">
          <span className="inline-flex items-center text-xs px-2 py-1 rounded-full border">
            <Eye className="w-3 h-3 mr-1" /> {project.visibility}
          </span>
        </div>
      </div>

      <div className="flex items-center gap-4 text-sm text-gray-600 mt-3">
        <span className="inline-flex items-center gap-1">
          <Package className="w-4 h-4" /> {project.icon_count} icons
        </span>
        <span className="inline-flex items-center gap-1">
          <Calendar className="w-4 h-4" /> {new Date(project.created_at).toLocaleDateString()}
        </span>
      </div>

      {(onEdit || onDelete) && (
        <div className="mt-4 flex gap-2">
          {onEdit && (
            <button onClick={() => onEdit(project)} className="px-3 py-2 text-sm rounded-md border inline-flex items-center gap-1">
              <Pencil className="w-4 h-4" /> Edit
            </button>
          )}
          {onDelete && (
            <button 
              onClick={() => onDelete(project)} 
              className="px-3 py-2 text-sm rounded-md border inline-flex items-center gap-1 text-red-600 border-red-200 hover:bg-red-50 hover:border-red-300 cursor-pointer transition-colors duration-200"
            >
              <Trash2 className="w-4 h-4" /> Delete
            </button>
          )}
          <Link to={`/manager/projects/${project.id}`} className="px-3 py-2 text-sm rounded-md border inline-flex items-center gap-1">
            <EyeIcon className="w-4 h-4" /> Details
          </Link>
        </div>
      )}
    </Card>
  );
}
