FROM jupyter/minimal-notebook

RUN conda install -y cassandra-driver pandas pep8 pylint

ADD test_data/ $HOME/work/test_data
ADD example.ipynb $HOME/work
ADD *.py $HOME/work/

RUN chown -R $USER:$USER $HOME
