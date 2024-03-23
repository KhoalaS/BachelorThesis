from pyecharts import options as opts

def img_opts() -> opts.ToolboxOpts:
    return opts.ToolboxOpts(is_show=True, feature=opts.ToolBoxFeatureOpts(save_as_image=opts.ToolBoxFeatureSaveAsImageOpts(background_color="white", pixel_ratio=8)))

rule_names = {
    "kTiny": "Tiny",
    "kVertDom": "VD",
    "kEdgeDom": "ED",
    "kSmall": "Small",
    "kTri": "Tri",
    "kExtTri": "ETri",
    "kApVertDom": "AVD",
    "kApDoubleVertDom": "ADVD",
    "kSmallEdgeDegTwo": "SED2",
    "kSmallEdgeDegTwo2": "SED2*",
    "kFallback": "F3"
}