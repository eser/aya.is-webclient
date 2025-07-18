import { Metadata } from "next";
import { CreateProfileForm } from "./_components/create-profile-form.tsx";

export const metadata: Metadata = {
  title: "Create Profile",
  description: "Create a new profile on aya.is",
};

export default function CreateProfilePage() {
  return (
    <div className="container mx-auto max-w-2xl px-4 py-8">
      <div className="space-y-6">
        <div className="text-center">
          <h1 className="text-3xl font-bold">Create Profile</h1>
          <p className="text-muted-foreground mt-2">
            Create your new profile. The profile will be created in the selected locale,
            and you can add additional locales later.
          </p>
        </div>

        <CreateProfileForm />
      </div>
    </div>
  );
}
