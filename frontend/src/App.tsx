import { useState, useEffect, useCallback } from "react";
import ImageCard from "./component/ImageCard";
import Header from "./component/Header";
import UploadModal from "./component/UploadModal";
import Filters from "./component/Filters";
import NoImagesMessage from "./component/NoImagesMessage";
import NewImagesNotification from "./component/NewImagesNotification";
import useImages from "./hooks/useImages";
import useTheme from "./hooks/useTheme";
import useWebSocket from "./hooks/useWebSocket";

function App() {
  const { images, filters, postImage, addImage, addFilter, removeFilter } =
    useImages();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const { isDark, toggleTheme } = useTheme();
  const [newImageCount, setNewImageCount] = useState(0);
  const [isAtTop, setIsAtTop] = useState(true);

  useEffect(() => {
    const handleScroll = () => {
      const atTop = window.scrollY < 100;
      setIsAtTop(atTop);
      if (atTop) {
        setNewImageCount(0);
      }
    };
    window.addEventListener("scroll", handleScroll, { passive: true });
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  const handleWebSocketMessage = useCallback(
    (image: (typeof images)[0]) => {
      addImage(image);
      if (!isAtTop) {
        setNewImageCount((prev) => prev + 1);
      }
    },
    [addImage, isAtTop],
  );

  useWebSocket({ onMessage: handleWebSocketMessage });

  const scrollToTop = () => {
    window.scrollTo({ top: 0, behavior: "smooth" });
    setNewImageCount(0);
  };

  return (
    <div className="max-w-lg w-full mx-auto min-h-screen">
      <Header
        onAddClick={() => setIsModalOpen(true)}
        isDark={isDark}
        onToggleTheme={toggleTheme}
      />

      <NewImagesNotification count={newImageCount} onClick={scrollToTop} />

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
        onSuccess={scrollToTop}
        onUpload={postImage}
      />
    </div>
  );
}

export default App;
