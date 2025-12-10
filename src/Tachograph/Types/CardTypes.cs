using System;

namespace Tachograph
{
    /// <summary>
    /// Represents the type of card.
    /// </summary>
    public enum CardType
    {
        Unknown = 0,
        DriverCard = 1,
        WorkshopCard = 2,
        ControlCard = 3,
        CompanyCard = 4
    }

    /// <summary>
    /// Represents a raw card file.
    /// </summary>
    public class RawCardFile
    {
        /// <summary>
        /// Gets or sets the raw binary data of the card file.
        /// </summary>
        public byte[] Data { get; set; }

        /// <summary>
        /// Creates a deep clone of this RawCardFile.
        /// </summary>
        /// <returns>A cloned RawCardFile</returns>
        public RawCardFile Clone()
        {
            return new RawCardFile
            {
                Data = (byte[])this.Data?.Clone()
            };
        }
    }

    /// <summary>
    /// Represents a parsed driver card file.
    /// </summary>
    public class DriverCardFile
    {
        // TODO: Add driver card specific fields
    }

    /// <summary>
    /// Represents a parsed workshop card file.
    /// </summary>
    public class WorkshopCardFile
    {
        // TODO: Add workshop card specific fields
    }

    /// <summary>
    /// Represents a parsed control card file.
    /// </summary>
    public class ControlCardFile
    {
        // TODO: Add control card specific fields
    }

    /// <summary>
    /// Represents a parsed company card file.
    /// </summary>
    public class CompanyCardFile
    {
        // TODO: Add company card specific fields
    }
}
