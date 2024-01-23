from pyecharts import options as opts

def img_opts() -> opts.ToolboxOpts:
    return opts.ToolboxOpts(is_show=True, feature=opts.ToolBoxFeatureOpts(save_as_image=opts.ToolBoxFeatureSaveAsImageOpts(background_color="white", pixel_ratio=2)))