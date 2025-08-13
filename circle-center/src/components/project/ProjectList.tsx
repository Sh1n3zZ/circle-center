import { useEffect, useState } from "react";
import ProjectListCard from "@/components/project/ProjectListCard";
import ProjectListCardCreate from "@/components/project/ProjectListCardCreate";
import { projectApi } from "@/api/manager/project";
import type { ProjectModel } from "@/api/manager/types";
import { toast } from "sonner";
import { LoadingSpinner } from "@/components/ui/loading-spinner";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";

export default function ProjectList() {
  const [projects, setProjects] = useState<ProjectModel[]>([]);
  const [loading, setLoading] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [projectToDelete, setProjectToDelete] = useState<ProjectModel | null>(null);

  const fetchData = async () => {
    try {
      setLoading(true);
      const res = await projectApi.listMyProjects();
      setProjects(res.data || []);
    } catch (e: any) {
      toast.error(e?.response?.data?.message || e.message || "Failed to load projects");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleDeleteClick = (project: ProjectModel) => {
    setProjectToDelete(project);
    setDeleteDialogOpen(true);
  };

  const handleDeleteConfirm = async () => {
    if (!projectToDelete) return;
    
    try {
      await projectApi.deleteProject(projectToDelete.id);
      toast.success("Project deleted");
      fetchData();
    } catch (e: any) {
      toast.error(e?.response?.data?.message || e.message || "Delete failed");
    } finally {
      setDeleteDialogOpen(false);
      setProjectToDelete(null);
    }
  };

  const handleDeleteCancel = () => {
    setDeleteDialogOpen(false);
    setProjectToDelete(null);
  };

  return (
    <>
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <div className="h-full">
          <ProjectListCardCreate onCreated={fetchData} />
        </div>
        {loading ? (
          <div className="col-span-full flex justify-center py-8">
            <LoadingSpinner className="text-gray-400" />
          </div>
        ) : projects.length === 0 ? (
          <div className="col-span-full text-center text-gray-500">No projects yet</div>
        ) : (
          projects.map((p) => (
            <ProjectListCard key={p.id} project={p} onDelete={handleDeleteClick} />
          ))
        )}
      </div>

      <AlertDialog open={deleteDialogOpen} onOpenChange={setDeleteDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Project</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete "{projectToDelete?.name}"? This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel onClick={handleDeleteCancel} className="cursor-pointer">Cancel</AlertDialogCancel>
            <AlertDialogAction onClick={handleDeleteConfirm} className="bg-red-600 hover:bg-red-700 cursor-pointer">
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}
