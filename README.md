# mmd-auto-trace-5

## ビルド

```bash
~/auto-trace/mmd-auto-trace-5/go/cmd$ go build -o mat5 main.go
~/auto-trace/mmd-auto-trace-5/go_sound/cmd$ go build -o sound main.go
```

## Human3R 環境構築 (WSL)

```bash
conda create -n human3r python=3.11 -y
conda activate human3r

python -m pip install --upgrade pip setuptools wheel

python -m pip uninstall -y numpy
python -m pip install numpy==1.26.4

python -m pip install torch==2.6.0 torchvision==0.21.0 torchaudio==2.6.0 --index-url https://download.pytorch.org/whl/cu126

python -m pip install -r requirements.txt
python -m pip install --no-build-isolation git+https://github.com/mattloper/chumpy@9b045ff5d6588a24a0bab52c83f032e2ba433e17

conda install -y 'llvm-openmp<16'

python -m pip install tqdm scikit-image gsplat gdown huggingface-hub

conda install -c conda-forge gcc=13 gxx=13

cd src/croco/models/curope/
python setup.py build_ext --inplace
cd ../../../../
```

```Bash
# SMPLX family models
bash scripts/fetch_smplx.sh

# Human3R checkpoints
hf download faneggg/human3r human3r_896L.pth --local-dir ./src
```

```bash
clear && CUDA_VISIBLE_DEVICES=0 python demo.py \
    --model_path src/human3r_896L.pth --size 512 \
    --seq_path ../sources/45seconds/45seconds_452-652.mp4 --subsample 1 --use_ttt3r \
    --vis_threshold 2 --downsample_factor 1 --reset_interval 100 \
    --save --output_dir ../sources/45seconds/output/ --render --render_video
```












## PromptHMR 環境構築 (WSL)

```bash
git clone https://github.com/miu200521358/mmd-PromptHMR.git
cd mmd-PromptHMR
```

```bash
conda create -n phmr python=3.12 -y
conda activate phmr

python -m pip install --upgrade pip setuptools wheel
python -m pip install numpy==1.26.4

python -m pip install torch==2.6.0 torchvision==0.21.0 torchaudio==2.6.0 --index-url https://download.pytorch.org/whl/cu126

python -m pip install --upgrade setuptools pip
python -m pip install torch-scatter -f https://data.pyg.org/whl/torch-2.6.0+cu126.html

python -m pip install -r requirements.txt
mkdir python_libs
git clone https://github.com/Arthur151/chumpy python_libs/chumpy
python -m pip install -e python_libs/chumpy --no-build-isolation

python -m pip install -U xformers==0.0.29.post2 --index-url https://download.pytorch.org/whl/cu126 --no-deps

gdown --folder -O ./data/ https://drive.google.com/drive/folders/151gPvMaUWok_pDQT6h8Rpvk_rCcKvcWZ?usp=sharing

python -m pip install data/wheels/sam2-1.6-cp312-cp312-linux_x86_64.whl
python -m pip install data/wheels/detectron2-0.9-cp312-cp312-linux_x86_64.whl
python -m pip install data/wheels/droid_backends_intr-0.4-cp312-cp312-linux_x86_64.whl
python -m pip install data/wheels/lietorch-0.4-cp312-cp312-linux_x86_64.whl
```

```bash
python scripts/demo_video.py --input_video data/examples/dance_1.mp4 --static_camera --viser_subsample 4
```













