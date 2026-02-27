import React from "react";
import { MainMenu } from "../../packages/excalidraw/index";
import { LanguageList } from "./LanguageList";

export const AppMainMenu: React.FC<{
  onCollabDialogOpen: () => any;
  isCollaborating: boolean;
  isCollabEnabled: boolean;
  onWorkspaceOpen: () => void;
  onWorkspaceSave: () => void;
  onWorkspaceSaveAs: () => void;
  onNewFile: () => void;
  workspaceFileName: string | null;
}> = React.memo((props) => {
  return (
    <MainMenu>
      <MainMenu.DefaultItems.LoadScene />
      <MainMenu.DefaultItems.SaveToActiveFile />
      <MainMenu.DefaultItems.Export />
      <MainMenu.DefaultItems.SaveAsImage />
      {props.isCollabEnabled && (
        <MainMenu.DefaultItems.LiveCollaborationTrigger
          isCollaborating={props.isCollaborating}
          onSelect={() => props.onCollabDialogOpen()}
        />
      )}
      <MainMenu.Separator />
      <MainMenu.Item onSelect={() => props.onNewFile()}>
        📄 New
      </MainMenu.Item>
      <MainMenu.Item onSelect={() => props.onWorkspaceOpen()}>
        📂 Open Workspace File…
      </MainMenu.Item>
      <MainMenu.Item onSelect={() => props.onWorkspaceSave()}>
        💾 {props.workspaceFileName
          ? `Save "${props.workspaceFileName}"`
          : "Save to Workspace…"}
      </MainMenu.Item>
      <MainMenu.Item onSelect={() => props.onWorkspaceSaveAs()}>
        📋 Save As…
      </MainMenu.Item>
      <MainMenu.Separator />
      <MainMenu.DefaultItems.Help />
      <MainMenu.DefaultItems.ClearCanvas />
      <MainMenu.DefaultItems.Socials />
      <MainMenu.Separator />
      <MainMenu.DefaultItems.ToggleTheme />
      <MainMenu.ItemCustom>
        <LanguageList style={{ width: "100%" }} />
      </MainMenu.ItemCustom>
      <MainMenu.DefaultItems.ChangeCanvasBackground />
    </MainMenu>
  );
});
