```yaml
instance:
    mtk instance create localdev
        - Create a new folder in working directory with argument
        - Add projects.localdev:`pwd` to ~/.mtk_config

    mtk instance add localdev
        - Add projects.localdev:`pwd` to ~/.mtk_config

    mtk instance switch localdev
        - Echo out `export $MTK_ROOT=`~/.mtk_config:projects.localdev`

environment:
    mtk environment switch development
        - Change environment in .mtk_config to argument

marketplace:
    mtk marketplace registry add git@github.com:msplat/msplat-marketplace.git
        - Add argument to $MTK_ROOT/.mtk_env_config:marketplace_registries

    mtk marketplace install bum_gateway
        - Get repo url from marketplaces
        - Git clone configuration@bum_gateway into $MTK_ROOT
        
projects:
    mtk projects clone
        - Git clone all repos from configuration@bum_gateway/environment/projects

    mtk projects build
    mtk projects list

services:
    mtk services restart
    mtk services logs

stacks:
    mtk stacks up
    mtk stacks rm
```