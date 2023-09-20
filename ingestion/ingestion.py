from langchain.document_loaders import WebBaseLoader
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain.embeddings import OpenAIEmbeddings
from langchain.vectorstores import Weaviate
import datetime, psycopg2, os
import warnings

STATUS_NEW="NEW"
STATUS_INGESTING="INGESTING"
STATUS_INGESTED="INGESTED"

class LaunchError(Exception):
    pass

def main():

    warnings.filterwarnings("ignore")

    envars = getEnv()

    conn = open_connection(envars["dbname"],envars["user"],envars["password"])
    cur = conn.cursor()
    
    resources = fetch_new_resources(cur)

    if resources is None or len(resources) == 0:
        print("No new resources found")
        return

    for resource in resources:

        print("Starting to ingest resource %s into domain %s" % (resource[1],resource[2]))

        mark_resource_as_ingesting(cur, conn, resource)

        loader = WebBaseLoader(resource[1])
        data = loader.load()

        text_splitter = RecursiveCharacterTextSplitter(chunk_size = 500, chunk_overlap = 0)
        all_splits = text_splitter.split_documents(data)
        Weaviate.from_documents(documents=all_splits, embedding=OpenAIEmbeddings(), index_name=resource[2], weaviate_url=envars["weaviate_host"])

        mark_resource_as_ingested(cur, conn, resource)

        print("Completed ingesting resource %s into domain %s" % (resource[1],resource[2]))

    close_connection(cur, conn)


def getEnv():

    envars = dict()

    envars["dbname"] = os.getenv("kh_ingestion_dbname", "")

    if envars["dbname"] == "":
        raise LaunchError("environment variable kh_ingestion_dbname not found")

    envars["user"] = os.getenv("kh_ingestion_user", "")

    if envars["user"] == "":
        raise LaunchError("environment variable kh_ingestion_user not found")
    
    envars["password"] = os.getenv("kh_ingestion_password", "")

    if envars["password"] == "":
        raise LaunchError("environment variable kh_ingestion_password not found")

    envars["weaviate_host"] = os.getenv("kh_ingestion_weaviate_host", "")

    if envars["weaviate_host"] == "":
        raise LaunchError("environment variable kh_ingestion_weaviate_host not found")

    envars["openai"] = os.getenv("kh_ingestion_openai_api_key", "")

    if envars["openai"] == "":
        raise LaunchError("environment variable kh_ingestion_openai_api_key not found")
    
    envars["OPENAI_API_KEY"] = envars["openai"]

    os.environ.update(envars)

    return envars

def open_connection(dbname, user, password):
    return psycopg2.connect("dbname=%s user=%s password=%s" % (dbname, user, password))

def close_connection(cur, conn):
    cur.close()
    conn.close()

def fetch_new_resources(cur):
    cur.execute("""SELECT id, url, domain_id FROM resources WHERE status = %s""", (STATUS_NEW,))
    records = cur.fetchall()
    return records

def now():
    datetime.datetime.now()

def mark_resource_as_ingesting(cur, conn, resource):
    timeNow = now()
    cur.execute("""UPDATE resources SET status=%s, ingestion_started_at=%s, updated_at=%s WHERE id=%s""", (STATUS_INGESTING, timeNow, timeNow, resource[0]))
    conn.commit()

def mark_resource_as_ingested(cur, conn, resource):
    timeNow = now()
    cur.execute("""UPDATE resources SET status=%s, ingestion_completed_at=%s, updated_at=%s WHERE id=%s""", (STATUS_INGESTED, timeNow, timeNow, resource[0]))
    conn.commit()


if __name__ == '__main__':
    main()