import {Button, Typography} from "@material-ui/core";
import {createStyles, makeStyles, Theme} from '@material-ui/core/styles';
import Toolbar from '@material-ui/core/Toolbar';
import React, {useState} from 'react';
import {resetServerContext} from "react-beautiful-dnd";
import {useSelector} from "react-redux";
import {indexActionCreators} from "../actions";
import {ImageListDrawer} from "../components/ImageListDrawer";
import {TagListDrawer} from "../components/TagListDrawer";
import {useActions, useVirtualizedAsset} from "../hooks";
import {Tag, WorkSpace} from "../models/models";
import {State} from "../reducers/reducer";
import {assetPathToUrl} from "../util";
import {tagActionCreators} from "../actions/tag";

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    content: {
      flexGrow: 1,
      padding: theme.spacing(3),
    },
    root: {
      display: 'flex',
    },
  }),
);

const useHandlers = (localState: LocalState, setLocalState: (s: LocalState) => void, globalState: GlobalState) => {
  const indexActionDispatcher = useActions(indexActionCreators);
  const tagActionDispatcher = useActions(tagActionCreators);
  return {
    ...indexActionDispatcher,
    ...tagActionDispatcher,
    clickAddDirectoryButton: (ws: WorkSpace) => {
      indexActionDispatcher.clickAddDirectoryButton({workSpaceName: ws.name});
    },
    clickAddTagButton: () => {
      const tag: Tag = {id: globalState.tags.length +1, name: ''};
      indexActionDispatcher.clickAddTagButton(tag);
      setLocalState({...localState, editTagId: tag.id});
    },
    renameTag: (tag: Tag) => {
      const workSpaceName = globalState.currentWorkSpace?.name!;
      if (workSpaceName === undefined) {
        // tslint:disable-next-line:no-console
        console.warn('workspace is null but tag is renamed', tag);
      }
      tagActionDispatcher.rename({workSpaceName, tag});
      setLocalState({...localState, editTagId: null});
    },
    clickEditTagButton: (tag: Tag) => {
      setLocalState({...localState, editTagId: tag.id});
    },
    clickImage: (selectedImgPath: string) => {
      setLocalState({...localState, selectedImgPath})
    },
    updateTags: (tags: Tag[]) => {
      const workSpaceName = globalState.currentWorkSpace?.name!;
      if (workSpaceName === undefined) {
        // tslint:disable-next-line:no-console
        console.warn('workspace is null but tags are updated', tags);
      }
      tagActionDispatcher.update({workSpaceName, tags});
    }
  };
}

interface GlobalState {
  currentWorkSpace: WorkSpace | null
  tags: Tag[]
  imagePaths: string[]
  isLoadingWorkSpace: boolean
  isScanningDirectories: boolean
}

const selector = (state: State): GlobalState => ({
  currentWorkSpace: state.global.currentWorkSpace,
  tags: state.global.tags,
  imagePaths: state.global.assets.map((a) => assetPathToUrl(a.path)),
  isLoadingWorkSpace: state.global.isLoadingWorkSpaces,
  isScanningDirectories: state.indexPage.scanning,
})

interface LocalState {
  // tags: Tag[]
  // maxId: number
  editTagId: number | null
  selectedImgPath: string | null
}

const generateInitialLocalState = (): LocalState => {
  return {
    editTagId: null,
    selectedImgPath: null,
  }
};

export default function Test() {
  const classes = useStyles();
  const [localState, setLocalState] = useState(generateInitialLocalState());
  const globalState = useSelector(selector)
  const handlers = useHandlers(localState, setLocalState, globalState);
  const virtualizedAssetProps = useVirtualizedAsset();
  const handleClickAddDirectoryButton = () => {
    if (globalState.currentWorkSpace === null) {
      // tslint:disable-next-line:no-console
      console.warn('workspace is not selected, but AddDirectoryButton is clicked')
      return
    }
    handlers.clickAddDirectoryButton(globalState.currentWorkSpace)
  }

  return (
    <div className={classes.root}>
      <ImageListDrawer
        {...virtualizedAssetProps}
        imagePaths={globalState.imagePaths}
        onClickImage={handlers.clickImage}
      />
      <main className={classes.content}>
        <Toolbar/>
        {localState.selectedImgPath === null ? null : <img
          src={localState.selectedImgPath}
          alt={localState.selectedImgPath}
          width='100%'
        />}
        <Typography paragraph={true}>
          Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt
          ut labore et dolore magna aliqua. Rhoncus dolor purus non enim praesent elementum
          facilisis leo vel. Risus at ultrices mi tempus imperdiet. Semper risus in hendrerit
          gravida rutrum quisque non tellus. Convallis convallis tellus id interdum velit laoreet id
          donec ultrices. Odio morbi quis commodo odio aenean sed adipiscing. Amet nisl suscipit
          adipiscing bibendum est ultricies integer quis. Cursus euismod quis viverra nibh cras.
          Metus vulputate eu scelerisque felis imperdiet proin fermentum leo. Mauris commodo quis
          imperdiet massa tincidunt. Cras tincidunt lobortis feugiat vivamus at augue. At augue eget
          arcu dictum varius duis at consectetur lorem. Velit sed ullamcorper morbi tincidunt. Lorem
          donec massa sapien faucibus et molestie ac.
        </Typography>
        <Button variant="outlined" color="primary"
                disabled={globalState.isScanningDirectories || globalState.isLoadingWorkSpace}
                onClick={handleClickAddDirectoryButton}>
          {globalState.isScanningDirectories ? 'Scanning...' : 'Add Directory'}
        </Button>
      </main>
      <TagListDrawer
        tags={globalState.tags}
        editTagId={localState.editTagId ?? undefined}
        onClickAddButton={handlers.clickAddTagButton}
        onClickEditButton={handlers.clickEditTagButton}
        onRename={handlers.renameTag}
        onUpdate={handlers.updateTags}
      />
    </div>
  );
}

export async function getServerSideProps() {
  resetServerContext()
  return {props: {}}
}
