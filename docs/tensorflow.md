# Install Tensorflow

### Make sure the correct version of Python3 and pip are installed.

## Install Tensorflow on Ubuntu-Focal

1. sudo apt update
2. sudo apt install -y python3-dev python3-pip python3-venv python3-testresources
3. python3 -m venv --system-site-packages ./tensorflow-env
4. source  ./tensorflow-env/bin/activate
5. pip install --upgrade pip
6. pip install numpy --upgrade
7. [Custom Install](https://github.com/lakshayg/tensorflow-build) pip install --ignore-installed --upgrade https://github.com/lakshayg/tensorflow-build/releases/download/tf2.4.0-ubuntu20.04-python3.8.5/tensorflow-2.4.1-cp38-cp38-linux_x86_64.whl --user
8. python -c "import tensorflow as tf;print(tf.reduce_sum(tf.random.normal([1000, 1000])))"


## Install Tensorflow on Ubuntu-Bionic

1. sudo apt update
2. sudo add-apt-repository ppa:deadsnakes/ppa
3. sudo apt install -y python3.7 python3-venv python3.7-venv python3-dev python3-testresources
4. python3.7 -m venv --system-site-packages ./tensorflow-env
5. source ./tensorflow-env/bin/activate
5. pip install --upgrade pip
6. pip install numpy --upgrade
7. [Custom Install](https://github.com/lakshayg/tensorflow-build) pip install --ignore-installed --upgrade https://github.com/lakshayg/tensorflow-build/releases/download/tf2.2.0-py3.7-ubuntu18.04/tensorflow-2.2.0-cp37-cp37m-linux_x86_64.whl --user
8. python -c "import tensorflow as tf;print(tf.reduce_sum(tf.random.normal([1000, 1000])))"

## Install Tensorflow on Darwin

1. brew reinstall python@3.9
2. python3.9 -m venv --system-site-packages ./tensorflow-env
3. source ./tensorflow-env/bin/activate
4. pip install --upgrade pip
5. pip install --upgrade tensorflow
6. python -c "import tensorflow as tf;print(tf.reduce_sum(tf.random.normal([1000, 1000])))"


