import { Navbar } from "@/components/navbar";
import { Footer } from "@/components/footer";
import { NotFoundView } from "@/components/not-found-view";

export default function NotFound() {
  return (
    <>
      <Navbar />
      <NotFoundView />
      <Footer />
    </>
  );
}
