import {CircularProgress} from "@nextui-org/progress";

export default function LoadingComponent() {
  return (
    <div
      style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        height: "100vh", // Make the loading spinner take up full viewport height
      }}
    >
      <CircularProgress />
    </div>
  );
}