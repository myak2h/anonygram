import { useState } from "react";
import ImageCard from "./component/ImageCard";
import Header from "./component/Header";
import UploadModal from "./component/UploadModal";
import Filters from "./component/Filters";
import NoImagesMessage from "./component/NoImagesMessage";
import useImages from "./hooks/useImages";
import useTheme from "./hooks/useTheme";

function App() {
  const { images, filters, postImage, addFilter, removeFilter } = useImages();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const { isDark, toggleTheme } = useTheme();

  return (
    <div className="max-w-lg w-full mx-auto min-h-screen">
      <Header
        onAddClick={() => setIsModalOpen(true)}
        isDark={isDark}
        onToggleTheme={toggleTheme}
      />

      <main className="flex flex-col gap-4 p-4 pt-20">
        <Filters
          filters={filters}
          onAddFilter={addFilter}
          onRemoveFilter={removeFilter}
        />

        {images.length === 0 ? (
          <NoImagesMessage />
        ) : (
          images.map((image) => (
            <div key={image.id} className="w-full">
              <ImageCard image={image} onTagClick={addFilter} />
            </div>
          ))
        )}
      </main>

      <UploadModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onSuccess={() => window.scrollTo({ top: 0, behavior: "smooth" })}
        onUpload={postImage}
      />
    </div>
  );
}

export default App;
