import GridList from '@material-ui/core/GridList';
import GridListTile from '@material-ui/core/GridListTile';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
import {VirtualizedAssetProps} from "../services/virtualizedAsset";
import {AssetListItemProps, VirtualizedAssetList} from "./VirtualizedAssetList";
import {assetPathToUrl} from "../util";

const useStyles = makeStyles((theme) => ({
  gridList: {
    height: '100%',
  },
  gridListTile: {
    cursor: 'pointer',
  },
  root: {
    backgroundColor: theme.palette.background.paper,
    display: 'flex',
    flexWrap: 'wrap',
    justifyContent: 'space-around',
    overflow: 'hidden',
  },
}));

interface Props extends VirtualizedAssetProps {
  cellHeight: number;
  width: number;
  paths: string[]
  onClickImage: (path: string) => void
}

// tslint:disable-next-line:variable-name
export const ImageGridList: React.FC<Props> = (props) => {
  const classes = useStyles();

  const genClickImageHandler = (imgPath: string) => () => {
    props.onClickImage(imgPath);
  };

  // tslint:disable-next-line:variable-name
  const ImageGridTile: React.FC<AssetListItemProps> = ({asset, isLoaded, style}) => {
    if (!isLoaded) {
      return (<div style={style}>Loading...</div>);
    }

    const pathUrl = assetPathToUrl(asset.path);

    return (
      <GridListTile
        style={style}
        key={asset.path}
        cols={1}
        onClick={genClickImageHandler(pathUrl)}
        className={classes.gridListTile}
      >
        <img src={pathUrl} />
      </GridListTile>
    );
  };

  return (
    <div className={classes.root}>
      <GridList cellHeight={props.cellHeight} className={classes.gridList} cols={1}>
        <VirtualizedAssetList
          assets={props.assets}
          hasMoreAssets={props.hasMoreAssets}
          isScanningAssets={props.isScanningAssets}
          onRequestNextPage={props.onRequestNextPage}
          workspace={props.workspace}
          height={window.innerHeight}
          itemSize={props.cellHeight}
          width={props.width}
        >
          {ImageGridTile}
        </VirtualizedAssetList>
      </GridList>
    </div>
  );
}