import { useState } from "react";
import { Plus, Loader2 } from "lucide-react";
import { projectApi } from "@/api/manager/project";
import type { CreateProjectRequest, ProjectVisibility } from "@/api/manager/types";
import { toast } from "sonner";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription, SheetTrigger, SheetFooter } from "@/components/ui/sheet";
import { Separator } from "@/components/ui/separator";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";

export default function ProjectListCardCreate({ onCreated }: { onCreated?: () => void }) {
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);

  const [name, setName] = useState("");
  const [slug, setSlug] = useState("");
  const [pkg, setPkg] = useState("");
  const [visibility, setVisibility] = useState<ProjectVisibility>("private");
  const [description, setDescription] = useState("");

  const reset = () => {
    setName("");
    setSlug("");
    setPkg("");
    setVisibility("private");
    setDescription("");
  };

  const handleSubmit = async () => {
    if (!name.trim()) {
      toast.error("Project name is required");
      return;
    }
    const payload: CreateProjectRequest = {
      name: name.trim(),
      slug: slug.trim() || undefined,
      package_name: pkg.trim() || undefined,
      visibility,
      description: description.trim() || undefined,
    };
    try {
      setLoading(true);
      await projectApi.createProject(payload);
      toast.success("Project created");
      setOpen(false);
      reset();
      onCreated?.();
    } catch (e: any) {
      toast.error(e?.response?.data?.message || e.message || "Create failed");
    } finally {
      setLoading(false);
    }
  };

  const isSubmitDisabled = loading || !name.trim();

  return (
    <>
      <Card className="h-full">
        <CardHeader className="gap-2">
          <CardTitle>Create project</CardTitle>
          <CardDescription>Start a new icon pack project.</CardDescription>
        </CardHeader>
        <CardContent>
          <Sheet open={open} onOpenChange={setOpen}>
            <SheetTrigger asChild>
              <Button variant="outline" className="w-full justify-center">
                <Plus className="w-4 h-4" /> New project
              </Button>
            </SheetTrigger>
            <SheetContent side="right">
              <SheetHeader>
                <SheetTitle>New project</SheetTitle>
                <SheetDescription>Fill in details to create a project.</SheetDescription>
              </SheetHeader>
              <div className="p-4 space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="name">
                    Name <span className="text-destructive">*</span>
                  </Label>
                  <Input
                    id="name"
                    placeholder="My Icon Pack"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    required
                    aria-required="true"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="slug">Slug</Label>
                  <Input id="slug" placeholder="my-icon-pack" value={slug} onChange={(e) => setSlug(e.target.value)} />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="pkg">Package name</Label>
                  <Input id="pkg" placeholder="com.example.icons" value={pkg} onChange={(e) => setPkg(e.target.value)} />
                </div>
                <div className="space-y-1">
                  <div className="flex items-center justify-between py-1">
                    <Label htmlFor="visibility">Public visibility</Label>
                    <Switch
                      id="visibility"
                      checked={visibility === "public"}
                      onCheckedChange={(checked) => setVisibility(checked ? "public" : "private")}
                    />
                  </div>
                  <p className="text-xs text-muted-foreground">If enabled, your project will be visible to others.</p>
                </div>
                <div className="space-y-2">
                  <Label htmlFor="desc">Description</Label>
                  <Textarea id="desc" rows={4} value={description} onChange={(e) => setDescription(e.target.value)} />
                </div>
                <Separator />
                <div className="flex justify-end gap-2">
                  <Button variant="outline" onClick={() => { setOpen(false); reset(); }} disabled={loading}>
                    Cancel
                  </Button>
                  <Button onClick={handleSubmit} disabled={isSubmitDisabled} aria-disabled={isSubmitDisabled}>
                    {loading && <Loader2 className="w-4 h-4 animate-spin" />} Create
                  </Button>
                </div>
              </div>
              <SheetFooter />
            </SheetContent>
          </Sheet>
        </CardContent>
      </Card>
    </>
  );
}
