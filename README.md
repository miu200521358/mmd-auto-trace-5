# mmd-auto-trace-5


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
    --seq_path ../sources/buster/buster_0-1700.mp4 --subsample 1 --use_ttt3r \
    --vis_threshold 2 --downsample_factor 1 --reset_interval 100 \
    --save --output_dir ../sources/buster/output/ --render --render_video
```











