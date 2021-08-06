#!/bin/bash
echo "Syncing git projects..."

initdir=$PWD
projdir=../../../

# landing page
if [ -d "${projdir}/landing-page" ]; then 
    echo "Landing page..."
    cd "${projdir}/landing-page"
    git pull
    cd $initdir
fi

# companion app
if [ -d "${projdir}/elabox-companion" ]; then 
    echo "Companion app..."
    cd "${projdir}/elabox-companion"
    git pull
    cd $initdir
fi